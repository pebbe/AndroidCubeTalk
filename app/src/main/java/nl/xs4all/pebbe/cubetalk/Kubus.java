package nl.xs4all.pebbe.cubetalk;

import android.content.Context;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.opengl.GLES20;
import android.opengl.GLUtils;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.FloatBuffer;

public class Kubus {

    private final static int VERTEX_ARRAY_SIZE1 = (1 * 6) * 3;
    private final static int VERTEX_ARRAY_SIZE2 = (5 * 6) * 3;
    private final static int COLOR_ARRAY_SIZE1 = (1 * 6) * 2;
    private final static int COLOR_ARRAY_SIZE2 = (5 * 6) * 2;

    private FloatBuffer coordsBuffer1;
    private FloatBuffer coordsBuffer2;
    private FloatBuffer colorsBuffer1;
    private FloatBuffer colorsBuffer2;
    private final int mProgram;
    private int mPositionHandle;
    private int mColorHandle;
    private int mCubeSizeHandle;

    private final String vertexShaderCode = "" +
            "uniform mat4 uMVPMatrix;" +
            "uniform vec3 uCubeSize;" +
            "attribute vec3 position;" +
            "attribute vec2 color;" +
            "uniform vec3 rgb;" +
            "varying vec2 col;" +
            "varying vec3 cl;" +
            "void main() {" +
            "    gl_Position = uMVPMatrix * vec4(uCubeSize*position, 1);" +
            "    col = color;" +
            "    cl = rgb;" +
            "}";

    private final String fragmentShaderCode = "" +
            "precision mediump float;" +
            "uniform sampler2D texture;" +
            "varying vec2 col;" +
            "varying vec3 cl;" +
            "void main() {" +
            "    gl_FragColor = vec4(cl, 1.0) * texture2D(texture, vec2(col[0], col[1]));" +
            "}";

    static final int COORDS_PER_VERTEX = 3;
    static float Coords1[] = new float[VERTEX_ARRAY_SIZE1];
    static float Coords2[] = new float[VERTEX_ARRAY_SIZE2];
    private final int coordStride = COORDS_PER_VERTEX * 4; // 4 bytes per float

    static final int COLORS_PER_VERTEX = 2;
    static float Colors1[] = new float[COLOR_ARRAY_SIZE1];
    static float Colors2[] = new float[COLOR_ARRAY_SIZE2];
   private final int colorStride = COLORS_PER_VERTEX * 4; // 4 bytes per float

    private int vertexCount1;
    private int vertexCount2;

    private void punt1(float x, float y, float z, float xi, float yi) {
        Coords1[COORDS_PER_VERTEX * vertexCount1 + 0] = x;
        Coords1[COORDS_PER_VERTEX * vertexCount1 + 1] = y;
        Coords1[COORDS_PER_VERTEX * vertexCount1 + 2] = z;
        Colors1[COLORS_PER_VERTEX * vertexCount1 + 0] = xi;
        Colors1[COLORS_PER_VERTEX * vertexCount1 + 1] = yi;
        vertexCount1++;
    }

    private void punt2(float x, float y, float z, float xi, float yi) {
        Coords2[COORDS_PER_VERTEX * vertexCount2 + 0] = x;
        Coords2[COORDS_PER_VERTEX * vertexCount2 + 1] = y;
        Coords2[COORDS_PER_VERTEX * vertexCount2 + 2] = z;
        Colors2[COLORS_PER_VERTEX * vertexCount2 + 0] = xi;
        Colors2[COLORS_PER_VERTEX * vertexCount2 + 1] = yi;
        vertexCount2++;
    }

    public Kubus() {
        vertexCount1 = 0;
        vertexCount2 = 0;

        // gezicht

        punt1(-1, 1, 1, 0, 0);
        punt1(-1, -1, 1, 0, 1);
        punt1(1, -1, 1, 1, 1);
        punt1(-1, 1, 1, 0, 0);
        punt1(1, -1, 1, 1, 1);
        punt1(1, 1, 1, 1, 0);

        // hoofd

        // rechts 1
        punt2(-1, 1, -1, .2f, 0);
        punt2(-1, -1, -1, .2f, 1);
        punt2(-1, -1, 1, 0, 1);
        punt2(-1, 1, -1, .2f, 0);
        punt2(-1, -1, 1, 0, 1);
        punt2(-1, 1, 1, 0, 0);

        // achter 2
        punt2(1, 1, -1, .2f, 0);
        punt2(1, -1, -1, .2f, 1);
        punt2(-1, -1, -1, .4f, 1);
        punt2(1, 1, -1, .2f, 0);
        punt2(-1, -1, -1, .4f, 1);
        punt2(-1, 1, -1, .4f, 0);

        // links 3
        punt2(1, 1, 1, .6f, 0);
        punt2(1, -1, 1, .6f, 1);
        punt2(1, -1, -1, .4f, 1);
        punt2(1, 1, 1, .6f, 0);
        punt2(1, -1, -1, .4f, 1);
        punt2(1, 1, -1, .4f, 0);

        // boven 4
        punt2(-1, 1, -1, .6f, 1);
        punt2(-1, 1, 1, .6f, 0);
        punt2(1, 1, 1, .8f, 0);
        punt2(-1, 1, -1, .6f, 1);
        punt2(1, 1, 1, .8f, 0);
        punt2(1, 1, -1, .8f, 1);

        // onder 5
        punt2(-1, -1, 1, .8f, 1);
        punt2(-1, -1, -1, .8f, 0);
        punt2(1, -1, -1, 1, 0);
        punt2(-1, -1, 1, .8f, 1);
        punt2(1, -1, -1, 1, 0);
        punt2(1, -1, 1, 1, 1);


        ByteBuffer b1 = ByteBuffer.allocateDirect(VERTEX_ARRAY_SIZE1 * 4);
        b1.order(ByteOrder.nativeOrder());
        coordsBuffer1 = b1.asFloatBuffer();
        coordsBuffer1.put(Coords1);
        coordsBuffer1.position(0);

        ByteBuffer b2 = ByteBuffer.allocateDirect(VERTEX_ARRAY_SIZE2 * 4);
        b2.order(ByteOrder.nativeOrder());
        coordsBuffer2 = b2.asFloatBuffer();
        coordsBuffer2.put(Coords2);
        coordsBuffer2.position(0);

        ByteBuffer b3 = ByteBuffer.allocateDirect(COLOR_ARRAY_SIZE1 * 4);
        b3.order(ByteOrder.nativeOrder());
        colorsBuffer1 = b3.asFloatBuffer();
        colorsBuffer1.put(Colors1);
        colorsBuffer1.position(0);

        ByteBuffer b4 = ByteBuffer.allocateDirect(COLOR_ARRAY_SIZE2 * 4);
        b4.order(ByteOrder.nativeOrder());
        colorsBuffer2 = b4.asFloatBuffer();
        colorsBuffer2.put(Colors2);
        colorsBuffer2.position(0);

        int vertexShader = Util.loadShader(
                GLES20.GL_VERTEX_SHADER, vertexShaderCode);
        int fragmentShader = Util.loadShader(
                GLES20.GL_FRAGMENT_SHADER, fragmentShaderCode);

        mProgram = GLES20.glCreateProgram();             // create empty OpenGL Program
        GLES20.glAttachShader(mProgram, vertexShader);   // add the vertex shader to program
        Util.checkGlError("glAttachShader vertexShader");
        GLES20.glAttachShader(mProgram, fragmentShader); // add the fragment shader to program
        Util.checkGlError("glAttachShader fragmentShader");
        GLES20.glLinkProgram(mProgram);                  // create OpenGL program executables
        Util.checkGlError("glLinkProgram");

    }

    public void draw(float[] mvpMatrix, float red, float green, float blue, float[] cubesize, int texturehead, int textureface) {
        drawPart(mvpMatrix, red, green, blue, cubesize, 1, textureface);
        drawPart(mvpMatrix, red, green, blue, cubesize, 2, texturehead);
    }

    private void drawPart(float[] mvpMatrix, float red, float green, float blue, float[] cubesize, int part, int texture) {

        // Add program to OpenGL environment
        GLES20.glUseProgram(mProgram);
        Util.checkGlError("glUseProgram");

        GLES20.glActiveTexture(GLES20.GL_TEXTURE0);
        Util.checkGlError("glActiveTexture");

        GLES20.glBindTexture(GLES20.GL_TEXTURE_2D, texture);
        Util.checkGlError("glBindTexture");

        // Set filtering
        GLES20.glTexParameteri(GLES20.GL_TEXTURE_2D, GLES20.GL_TEXTURE_MIN_FILTER, GLES20.GL_LINEAR);
        Util.checkGlError("glTexParameteri");

        GLES20.glTexParameteri(GLES20.GL_TEXTURE_2D, GLES20.GL_TEXTURE_MAG_FILTER, GLES20.GL_LINEAR);
        Util.checkGlError("glTexParameteri");

        GLES20.glDisable(GLES20.GL_BLEND);

        mPositionHandle = GLES20.glGetAttribLocation(mProgram, "position");
        Util.checkGlError("glGetAttribLocation position");
        GLES20.glEnableVertexAttribArray(mPositionHandle);
        Util.checkGlError("glEnableVertexAttribArray position");
        GLES20.glVertexAttribPointer(
                mPositionHandle, COORDS_PER_VERTEX,
                GLES20.GL_FLOAT, false,
                coordStride, part == 1 ? coordsBuffer1 : coordsBuffer2);
        Util.checkGlError("glVertexAttribPointer position");

        mColorHandle = GLES20.glGetAttribLocation(mProgram, "color");
        Util.checkGlError("glGetAttribLocation color");
        GLES20.glEnableVertexAttribArray(mColorHandle);
        Util.checkGlError("glEnableVertexAttribArray color");
        GLES20.glVertexAttribPointer(
                mColorHandle, COLORS_PER_VERTEX,
                GLES20.GL_FLOAT, false,
                colorStride, part == 1 ? colorsBuffer1 : colorsBuffer2);
        Util.checkGlError("glVertexAttribPointer color");

        int mMatrixHandle = GLES20.glGetUniformLocation(mProgram, "uMVPMatrix");
        Util.checkGlError("glGetUniformLocation uMVPMatrix");
        GLES20.glUniformMatrix4fv(mMatrixHandle, 1, false, mvpMatrix, 0);
        Util.checkGlError("glUniformMatrix4fv uMVPMatrix");

        int mRgbHandle = GLES20.glGetUniformLocation(mProgram, "rgb");
        Util.checkGlError("glGetUniformLocation rgb");
        GLES20.glUniform3f(mRgbHandle, red, green, blue);
        Util.checkGlError("glUniformMatrix4fv rgb");

        mCubeSizeHandle = GLES20.glGetUniformLocation(mProgram, "uCubeSize");
        Util.checkGlError("glGetUniformLocation uCubeSize");
        GLES20.glUniform3f(mCubeSizeHandle, cubesize[0], cubesize[1], cubesize[2]);
        Util.checkGlError("glUniform3f uCubeSize");

        // Get handle to textures locations
        int mSamplerLoc = GLES20.glGetUniformLocation(mProgram, "texture");
        Util.checkGlError("glGetUniformLocation texture");
        // Set the sampler texture unit to 0, where we have saved the texture.
        GLES20.glUniform1i(mSamplerLoc, 0);
        Util.checkGlError("glUniform1i mSamplerLoc");

        // Draw
        GLES20.glDrawArrays(GLES20.GL_TRIANGLES, 0, part == 1 ? vertexCount1 : vertexCount2);
        Util.checkGlError("glDrawArrays");

        // Disable vertex arrays
        GLES20.glDisableVertexAttribArray(mColorHandle);
        Util.checkGlError("glDisableVertexAttribArray colorHandle");
        GLES20.glDisableVertexAttribArray(mPositionHandle);
        Util.checkGlError("glDisableVertexAttribArray positionHandle");
    }

}
