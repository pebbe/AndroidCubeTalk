package nl.xs4all.pebbe.cubetalk;

public class CubeData {
    public float[] pos;
    public float[] lookat;
    public float roll = 0;
    public float[] color;
    public boolean visible = false;
    public long idx_enter_exit = 0;
    public long idx_moveto = 0;
    public long idx_lookat = 0;
    public long idx_color = 0;
    public long idx_head = 0;
    public long idx_face = 0;
    public int texturehead;
    public int textureface;

    public CubeData(int head, int face) {
        pos = new float[]{0, 0, 100};
        lookat = new float[]{0, 0, -1};
        color = new float[]{1, 1, 1};
        texturehead = head;
        textureface = face;
    }
}
