package nl.xs4all.pebbe.cubetalk;

public class CubeData {
    public float[] pos;
    public float[] lookat;
    public float[] color;
    public boolean visible = false;
    public long idx_enter_exit = 0;
    public long idx_moveto = 0;
    public long idx_lookat = 0;
    public long idx_color = 0;
    public boolean valid;

    public CubeData() {
        pos = new float[]{0, 0, 100};
        lookat = new float[]{0, 0, -1};
        color = new float[]{1, 1, 1};
    }
}
