package nl.xs4all.pebbe.cubetalk;

import android.opengl.GLES20;
import android.util.Log;

public class Util {

    public final static int TEXTURE_WORLD = 0;
    public final static int TEXTURE_INFO0 = 1;
    public final static int TEXTURE_INFO1 = 2;
    public final static int TEXTURE_HEAD0 = 3;
    public final static int TEXTURE_HEAD1 = 4;
    public final static int TEXTURE_HEAD2 = 5;
    public final static int TEXTURE_HEAD3 = 6;
    public final static int TEXTURE_FACE0 = 7;
    public final static int TEXTURE_FACE1 = 8;
    public final static int TEXTURE_FACE2 = 9;
    public final static int TEXTURE_FACE3 = 10;

    public final static int NR_OF_TEXTURES = 11;


    public final static String kAddress = "address";
    public final static String kPort = "port";
    public final static String kUid = "uID";
    public final static String sError = "error";

    public final static int stOK = 0;
    public final static int stNIL = 1;
    public final static int stERROR = 2;

    public static int loadShader(int type, String shaderCode) {

        // create a vertex shader type (GLES20.GL_VERTEX_SHADER)
        // or a fragment shader type (GLES20.GL_FRAGMENT_SHADER)
        int shader = GLES20.glCreateShader(type);

        // add the source code to the shader and compile it
        GLES20.glShaderSource(shader, shaderCode);
        checkGlError("glShaderSource");
        GLES20.glCompileShader(shader);
        checkGlError("glCompileShader");

        return shader;
    }

    public static void checkGlError(String glOperation) {
        int error;
        while ((error = GLES20.glGetError()) != GLES20.GL_NO_ERROR) {
            Log.e("GL-ERROR", glOperation + ": glError " + error);
            throw new RuntimeException(glOperation + ": glError " + error);
        }
    }
}
