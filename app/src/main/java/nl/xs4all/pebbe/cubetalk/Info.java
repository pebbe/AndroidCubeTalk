package nl.xs4all.pebbe.cubetalk;

import android.content.Context;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.opengl.GLES20;
import android.opengl.GLUtils;
import android.util.Log;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.FloatBuffer;

import static java.lang.Math.PI;

public class Info {

    private final static int ARRAY_SIZE = 2 * 6 * 6;

    private FloatBuffer vertexBuffer;
    private final int mProgram;
    private int mPositionHandle;
    private int mMatrixHandle;
    private int texture;

    private final String vertexShaderCode = "" +
            "uniform mat4 uMVPMatrix;" +
            "attribute vec2 position;" +
            "varying vec2 color;" +
            "void main() {" +
            "    gl_Position = uMVPMatrix * vec4(position[0] * 2.0 - 1.0, (position[1] * 2.0 - 1.0) * 0.75, 0.0, 1.0);" +
            "    color = vec2(position[0], 1.0 - position[1]);" +
            "}";

    private final String fragmentShaderCode = "" +
            "precision mediump float;" +
            "uniform sampler2D texture;" +
            "varying vec2 color;" +
            "void main() {" +
            "    gl_FragColor = texture2D(texture, color);" +
            "}";

    static final int COORDS_PER_VERTEX = 2;
    static float Coords[] = new float[ARRAY_SIZE];
    private final int vertexStride = COORDS_PER_VERTEX * 4; // 4 bytes per vertex
    private int vertexCount;

    private void Point (float x, float y) {
        Coords[COORDS_PER_VERTEX * vertexCount + 0] = x;
        Coords[COORDS_PER_VERTEX * vertexCount + 1] = y;
        vertexCount++;
    }


    public Info(Context context, int texturename, String[] lines) {
        texture = texturename;

        vertexCount = 0;
        Point(0, 1);
        Point(0, 0);
        Point(1, 0);
        Point(0, 1);
        Point(1, 0);
        Point(1, 1);

        ByteBuffer bb = ByteBuffer.allocateDirect(ARRAY_SIZE * 4);
        bb.order(ByteOrder.nativeOrder());
        vertexBuffer = bb.asFloatBuffer();
        vertexBuffer.put(Coords);
        vertexBuffer.position(0);

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

        // Temporary create a bitmap
        Bitmap bmp = BitmapFactory.decodeResource(context.getResources(), R.raw.info);

        // Bind texture to texturename
        GLES20.glActiveTexture(GLES20.GL_TEXTURE0);
        Util.checkGlError("glActiveTexture");
        GLES20.glBindTexture(GLES20.GL_TEXTURE_2D, texture);
        Util.checkGlError("glBindTexture");

        // Load the bitmap into the bound texture.
        GLUtils.texImage2D(GLES20.GL_TEXTURE_2D, 0, bmp, 0);

        // We are done using the bitmap so we should recycle it.
        bmp.recycle();
    }

    public void draw(float[] mvpMatrix) {
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

        mPositionHandle = GLES20.glGetAttribLocation(mProgram, "position");
        Util.checkGlError("glGetAttribLocation vPosition");
        GLES20.glEnableVertexAttribArray(mPositionHandle);
        Util.checkGlError("glEnableVertexAttribArray position");
        GLES20.glVertexAttribPointer(
                mPositionHandle, COORDS_PER_VERTEX,
                GLES20.GL_FLOAT, false,
                vertexStride, vertexBuffer);
        Util.checkGlError("glVertexAttribPointer position");

        mMatrixHandle = GLES20.glGetUniformLocation(mProgram, "uMVPMatrix");
        Util.checkGlError("glGetUniformLocation uMVPMatrix");
        GLES20.glUniformMatrix4fv(mMatrixHandle, 1, false, mvpMatrix, 0);
        Util.checkGlError("glUniformMatrix4fv uMVPMatrix");

        // Get handle to textures locations
        int mSamplerLoc = GLES20.glGetUniformLocation (mProgram, "texture" );
        Util.checkGlError("glGetUniformLocation texture");
        // Set the sampler texture unit to 0, where we have saved the texture.
        GLES20.glUniform1i(mSamplerLoc, 0);
        Util.checkGlError("glUniform1i mSamplerLoc");

        // Draw
        GLES20.glDrawArrays(GLES20.GL_TRIANGLES, 0, vertexCount);
        Util.checkGlError("glDrawArrays");

        // Disable vertex arrays
        GLES20.glDisableVertexAttribArray(mPositionHandle);
        Util.checkGlError("glDisableVertexAttribArray mPositionHandle");
    }
}
