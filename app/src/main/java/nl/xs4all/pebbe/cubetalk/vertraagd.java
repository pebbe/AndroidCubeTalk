package nl.xs4all.pebbe.cubetalk;

import java.util.concurrent.BlockingQueue;
import java.util.concurrent.DelayQueue;

public class vertraagd implements MainActivity.Provider {

    private long delay = 1000;
    private float enhance = 1.0f;
    private BlockingQueue queue;

    public vertraagd(long d, float e) {
        delay = d;
        enhance = e;
        queue = new DelayQueue();
    }

    @Override
    public String getError() {
        return "";
    }

    @Override
    public int forward(float[] out, float[] in) {
        DelayObject object = new DelayObject(in[0], in[1], in[2], delay);
        try {
            queue.put(object);
        } catch (Exception ex) {
        }

        object = (DelayObject) queue.poll();
        if (object == null) {
            return Util.stNIL;
        }

        DelayObject o;
        while ((o = (DelayObject) queue.poll()) != null) {
            object = o;
        }

        if (enhance == 1.0f) {
            out[0] = object.x;
            out[1] = object.y;
            out[2] = object.z;
        } else {
            float roth = enhance * (float) Math.atan2((double) object.x, (double) -object.z);
            float rotv = /* enhance * */ (float) Math.atan2(object.y, Math.sqrt(object.x * object.x + object.z * object.z));
            double cosrotv =  Math.cos(rotv);
            out[0] = (float) (Math.sin(roth) * cosrotv);
            out[1] = (float) Math.sin(rotv);
            out[2] = (float) (-Math.cos(roth) * cosrotv);
        }
        return Util.stOK;
    }
}
