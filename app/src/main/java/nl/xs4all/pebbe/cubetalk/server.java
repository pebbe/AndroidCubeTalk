package nl.xs4all.pebbe.cubetalk;

import android.content.Context;

import java.io.DataInputStream;
import java.io.PrintStream;
import java.net.Socket;
import java.util.Locale;

public class server implements MainActivity.Provider {

    private static final int NR_OF_CONNECTIONS = 8;
    private int current = 0;

    private Socket[] sockets;
    private DataInputStream[] inputs;
    private PrintStream[] outputs;

    private float x;
    private float y;
    private float z;
    private boolean ok = false;
    private boolean err = false;
    private String ErrStr = "";
    final private Object xyzokLock = new Object();

    private boolean[] runnings;
    final private Object runningLock = new Object();

    public server(Context context, String address, int port) {
        sockets = new Socket[NR_OF_CONNECTIONS];
        inputs = new DataInputStream[NR_OF_CONNECTIONS];
        outputs = new PrintStream[NR_OF_CONNECTIONS];
        runnings = new boolean[NR_OF_CONNECTIONS];
        for (int i = 0; i < NR_OF_CONNECTIONS; i++) {
            runnings[i] = true;
        }

        MyDBHandler handler = new MyDBHandler(context);
        String value = handler.findSetting(Util.kUid);
        if (value.equals("")) {
            value = "" + System.currentTimeMillis();
            handler.addSetting(Util.kUid, value);
        }
        final String uid = value;
        final String addr = address;
        final int pnum = port;

        Runnable runnable = new Runnable() {
            @Override
            public void run() {
                for (int i = 0; i < NR_OF_CONNECTIONS; i++) {
                    try {
                        sockets[i] = new Socket(addr, pnum);
                        sockets[i].setSoTimeout(1000);
                        inputs[i] = new DataInputStream(sockets[i].getInputStream());
                        outputs[i] = new PrintStream(sockets[i].getOutputStream());
                        outputs[i].format("VRC1.0 %s\n", uid);
                        inputs[i].readLine();
                    } catch (Exception e) {
                        e.printStackTrace();
                    }
                    synchronized (runningLock) {
                        runnings[i] = false;
                    }
                }
            }
        };
        Thread thread = new Thread(runnable);
        thread.start();
    }

    @Override
    public String getError() {
        return ErrStr;
    }

    @Override
    public int forward(float[] out, float[] in) {
        final float xi = in[0];
        final float yi = in[1];
        final float zi = in[2];
        final int index = current;
        current = (current + 1) % NR_OF_CONNECTIONS;
        Runnable runnable = new Runnable() {
            @Override
            public void run() {
                synchronized (runningLock) {
                    if (runnings[index]) {
                        return;
                    }
                    runnings[index] = true;
                }
                outputs[index].format(Locale.US, "%f %f %f\n", xi, yi, zi);

                String response;
                try {
                    response = inputs[index].readLine();
                } catch (Exception e) {
                    synchronized (runningLock) {
                        runnings[index] = false;
                    }
                    synchronized (xyzokLock) {
                        err = true;
                        ErrStr = e.toString();
                    }
                    return;
                }

                synchronized (xyzokLock) {
                    ok = false;
                    if (response == null) {
                        err = true;
                        ErrStr = "No response from remote server";
                    } else {
                        String[] parts = response.trim().split("[ \t]+");
                        if (parts.length == 3) {
                            try {
                                x = Float.parseFloat(parts[0]);
                                y = Float.parseFloat(parts[1]);
                                z = Float.parseFloat(parts[2]);
                                ok = true;
                            } catch (Exception e) {
                                // shouldn't happen
                                err = true;
                                ErrStr = e.toString();
                            }
                        }
                    }
                }
                synchronized (runningLock) {
                    runnings[index] = false;
                }
            }
        };
        Thread thread = new Thread(runnable);
        thread.start();


        int retval = Util.stNIL;
        synchronized (xyzokLock) {
            if (err) {
                retval = Util.stERROR;
                err = false;
            } else if (ok) {
                out[0] = x;
                out[1] = y;
                out[2] = z;
                retval = Util.stOK;
            }
            ok = false;
        }
        return retval;
    }
}
