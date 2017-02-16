package nl.xs4all.pebbe.cubetalk;

import android.opengl.GLES20;
import android.util.Log;

public class Util {

    public final static String kAddress = "address";
    public final static String kPort = "port";
    public final static String kUid = "uid";
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
