package nl.xs4all.pebbe.cubetalk;

import java.util.concurrent.Delayed;
import java.util.concurrent.TimeUnit;

public class DelayObject implements Delayed {

    public float x;
    public float y;
    public float z;
    private long startTime;

    public DelayObject(float x, float y, float z, long delay) {
        this.x = x;
        this.y = y;
        this.z = z;
        this.startTime = System.currentTimeMillis() + delay;
    }

    @Override
    public long getDelay(TimeUnit unit) {
        long left = startTime - System.currentTimeMillis();
        return unit.convert(left, TimeUnit.MILLISECONDS);
    }

    @Override
    public int compareTo(Delayed o) {
        if (this.startTime < ((DelayObject) o).startTime) {
            return -1;
        }
        if (this.startTime > ((DelayObject) o).startTime) {
            return 1;
        }
        return 0;
    }
}
