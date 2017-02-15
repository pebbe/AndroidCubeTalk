package nl.xs4all.pebbe.cubetalk;

import android.content.Context;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.opengl.GLES20;
import android.opengl.GLUtils;

import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.nio.FloatBuffer;

public class Arrows {

    private final static int ARRAY_SIZE = 6 * 3 * 2 * 4;

    private FloatBuffer vertexBuffer;
    private final int mProgram;
    private int mPositionHandle;
    private int mMatrixHandle;

    private final String vertexShaderCode = "" +
            "uniform mat4 uMVPMatrix;" +
            "attribute vec4 position;" +
            "varying vec3 color;" +
            "void main() {" +
            "    gl_Position = uMVPMatrix * vec4(position[0], position[1], position[2], 1.0);" +
            "    color = vec3(position[3], position[3], 0.5 + 0.5 * position[3]);" +
            "}";

    private final String fragmentShaderCode = "" +
            "precision mediump float;" +
            "uniform sampler2D texture;" +
            "varying vec3 color;" +
            "void main() {" +
            "    gl_FragColor = vec4(color, 1.0);" +
            "}";

    static final int COORDS_PER_VERTEX = 4;
    static float Coords[] = new float[ARRAY_SIZE];
    private final int vertexStride = COORDS_PER_VERTEX * 4; // 4 bytes per vertex
    private int vertexCount;

    private void point(float x, float y, float r, float d, float c) {
        Coords[COORDS_PER_VERTEX * vertexCount] = -d * (float) Math.sin(r) + (float) Math.cos(-r) * x;
        Coords[COORDS_PER_VERTEX * vertexCount + 1] = y;
        Coords[COORDS_PER_VERTEX * vertexCount + 2] = -d * (float) Math.cos(r) + (float) Math.sin(-r) * x;
        Coords[COORDS_PER_VERTEX * vertexCount + 3] = c;
        vertexCount++;
   }

    public Arrows() {

        vertexCount = 0;
        for (int i = 0; i < 3; i++) {
            float r = (67.5f + 45f * (float)i) / 180.0f * (float)Math.PI;
            point(-.16f, .32f, r, 3.01f, 0);
            point(-.16f, -.32f, r, 3.01f, 0);
            point(.17f, 0, r, 3.01f, 0);

            point(.16f, .32f, -r, 3.01f, 0);
            point(-.17f, 0, -r, 3.01f, 0);
            point(.16f, -.32f, -r, 3.01f, 0);

            point(-.15f, .3f, r, 3.0f, 1);
            point(-.15f, -.3f, r, 3.0f, 1);
            point(.15f, 0, r, 3.0f, 1);

            point(.15f, .3f, -r, 3.0f, 1);
            point(-.15f, 0, -r, 3.0f, 1);
            point(.15f, -.3f, -r, 3.0f, 1);
        }

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

    }

    public void draw(float[] mvpMatrix) {
        // Add program to OpenGL environment
        GLES20.glUseProgram(mProgram);
        Util.checkGlError("glUseProgram");

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
