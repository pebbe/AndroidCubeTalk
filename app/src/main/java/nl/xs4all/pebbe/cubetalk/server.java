package nl.xs4all.pebbe.cubetalk;

import android.content.Context;

import java.io.DataInputStream;
import java.io.PrintStream;
import java.net.Socket;
import java.util.Locale;

public class server implements MainActivity.Provider {

    private static final int NR_OF_CONNECTIONS = 20;
    private int current = 0;

    private static final int MAX_CUBES = 100;
    private int nr_of_cubes = 0;

    private Socket[] sockets;
    private DataInputStream[] inputs;
    private PrintStream[] outputs;

    private String[] ids;
    private CubeData[] cubes;
    private float SelfPos = 4;
    private long SelfIdx = 0;
    private boolean err = false;
    private String ErrStr = "";
    final private Object settingsLock = new Object();

    private boolean[] runnings;
    final private Object runningLock = new Object();

    public server(Context context, String address, int port) {
        ids = new String[MAX_CUBES];
        cubes = new CubeData[MAX_CUBES];

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
                        outputs[i].format("join %s\n", uid);
                        inputs[i].readLine(); // .
                        if (i == 0) {
                            outputs[i].format("reset\n");
                            inputs[i].readLine(); // .
                        }
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
        synchronized (settingsLock) {
            return ErrStr;
        }
    }

    @Override
    public float getSelf() {
        synchronized (settingsLock) {
            return SelfPos;
        }
    }

    @Override
    public CubeData getCubeData(int i) {
        CubeData cube = new CubeData();
        synchronized (settingsLock) {
            if (i >= nr_of_cubes) {
                cube.valid = false;
            } else {
                cube.valid = true;
                cube.pos[0] = cubes[i].pos[0];
                cube.pos[1] = cubes[i].pos[1];
                cube.pos[2] = cubes[i].pos[2];
                cube.lookat[0] = cubes[i].lookat[0];
                cube.lookat[1] = cubes[i].lookat[1];
                cube.lookat[2] = cubes[i].lookat[2];
                cube.color[0] = cubes[i].color[0];
                cube.color[1] = cubes[i].color[1];
                cube.color[2] = cubes[i].color[2];
                cube.visible = cubes[i].visible;
            }
        }
        return cube;
    }

    @Override
    public int forward(float[] in) {
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
                outputs[index].format(Locale.US, "lookat %f %f %f\n", xi, yi, zi);

                boolean busy = true;
                while (busy) {

                    String response;
                    try {
                        response = inputs[index].readLine();
                    } catch (Exception e) {
                        synchronized (runningLock) {
                            runnings[index] = false;
                        }
                        synchronized (settingsLock) {
                            err = true;
                            ErrStr = e.toString();
                        }
                        return;
                    }

                    if (response == null) {
                        err = true;
                        ErrStr = "No response from remote server";
                        busy = false;
                        break;
                    }

                    String e = "";
                    String[] parts = response.trim().split("[ \t]+");
                    if (parts.length > 0) {
                        if (parts[0].equals(".")) {
                            busy = false;
                        } else if (parts[0].equals("lookat")) {
                            e = setLookat(parts);
                        } else if (parts[0].equals("self")) {
                            e = setSelf(parts);
                        } else if (parts[0].equals("enter")) {
                            e = setEnter(parts);
                        } else if (parts[0].equals("exit")) {
                            e = setExit(parts);
                        } else if (parts[0].equals("moveto")) {
                            e = setMoveto(parts);
                        } else if (parts[0].equals("color")) {
                            e = setColor(parts);
                        }
                    }
                    if (!e.equals("")) {
                        err = true;
                        ErrStr = e;
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
        synchronized (settingsLock) {
            if (err) {
                retval = Util.stERROR;
                err = false;
            }
        }
        return retval;
    }

    // self {n0} {z}
    private String setSelf(String[] parts) {
        if (parts.length == 3) {
            long n = 0;
            float z = 0;
            try {
                n = Integer.parseInt(parts[1], 10);
                z = Float.parseFloat(parts[2]);
            } catch (Exception e) {
                return e.toString();
            }
            synchronized (settingsLock) {
                if (n >= SelfIdx) {
                    SelfIdx = n;
                    SelfPos = z;
                }
            }
        }
        return "";
    }

    // enter {id} {n1}
    private String setEnter(String[] parts) {
        if (parts.length == 3) {
            synchronized (settingsLock) {
                long n = 0;
                try {
                    n = Integer.parseInt(parts[2]);
                } catch (Exception e) {
                    return e.toString();
                }
                boolean found = false;
                int i;
                for (i = 0; i < nr_of_cubes; i++) {
                    if (ids[i].equals(parts[1])) {
                        found = true;
                        break;
                    }
                }
                if (found) {
                    if (n >= cubes[i].idx_enter_exit) {
                        cubes[i].idx_enter_exit = n;
                        cubes[i].visible = true;
                    }
                } else {
                    if (nr_of_cubes < MAX_CUBES) {
                        cubes[nr_of_cubes] = new CubeData();
                        cubes[nr_of_cubes].visible = true;
                        ids[i] = parts[1];
                        nr_of_cubes++;
                    }
                }
            }
        }
        return "";
    }

    // exit {id} {n1}
    private String setExit(String[] parts) {
        if (parts.length == 3) {
            synchronized (settingsLock) {
                long n = 0;
                try {
                    n = Integer.parseInt(parts[2]);
                } catch (Exception e) {
                    return e.toString();
                }
                boolean found = false;
                int i;
                for (i = 0; i < nr_of_cubes; i++) {
                    if (ids[i].equals(parts[1])) {
                        found = true;
                        break;
                    }
                }
                if (found) {
                    if (n >= cubes[i].idx_enter_exit) {
                        cubes[i].idx_enter_exit = n;
                        cubes[i].visible = false;
                    }
                }
            }
        }
        return "";
    }

    // moveto {id} {n2} {x} {y} {z}
    private String setMoveto(String[] parts) {
        if (parts.length == 6) {
            synchronized (settingsLock) {
                long n = 0;
                float x;
                float y;
                float z;
                try {
                    n = Integer.parseInt(parts[2]);
                    x = Float.parseFloat(parts[3]);
                    y = Float.parseFloat(parts[4]);
                    z = Float.parseFloat(parts[5]);
                } catch (Exception e) {
                    return e.toString();
                }
                boolean found = false;
                int i;
                for (i = 0; i < nr_of_cubes; i++) {
                    if (ids[i].equals(parts[1])) {
                        found = true;
                        break;
                    }
                }
                if (found) {
                    if (n >= cubes[i].idx_moveto) {
                        cubes[i].idx_moveto = n;
                        cubes[i].pos[0] = x;
                        cubes[i].pos[1] = y;
                        cubes[i].pos[2] = z;
                    }
                }
            }
        }
        return "";
    }

    // lookat {id} {n3} {x} {y} {z}
    private String setLookat(String[] parts) {
        if (parts.length == 6) {
            synchronized (settingsLock) {
                long n = 0;
                float x;
                float y;
                float z;
                try {
                    n = Integer.parseInt(parts[2]);
                    x = Float.parseFloat(parts[3]);
                    y = Float.parseFloat(parts[4]);
                    z = Float.parseFloat(parts[5]);
                } catch (Exception e) {
                    return e.toString();
                }
                boolean found = false;
                int i;
                for (i = 0; i < nr_of_cubes; i++) {
                    if (ids[i].equals(parts[1])) {
                        found = true;
                        break;
                    }
                }
                if (found) {
                    if (n >= cubes[i].idx_lookat) {
                        cubes[i].idx_lookat = n;
                        cubes[i].lookat[0] = x;
                        cubes[i].lookat[1] = y;
                        cubes[i].lookat[2] = z;
                    }
                }
            }
        }
        return "";
    }

    // color {id} {n4} {red} {green} {blue}
    private String setColor(String[] parts) {
        if (parts.length == 6) {
            synchronized (settingsLock) {
                long n = 0;
                float r;
                float g;
                float b;
                try {
                    n = Integer.parseInt(parts[2]);
                    r = Float.parseFloat(parts[3]);
                    g = Float.parseFloat(parts[4]);
                    b = Float.parseFloat(parts[5]);
                } catch (Exception e) {
                    return e.toString();
                }
                boolean found = false;
                int i;
                for (i = 0; i < nr_of_cubes; i++) {
                    if (ids[i].equals(parts[1])) {
                        found = true;
                        break;
                    }
                }
                if (found) {
                    if (n >= cubes[i].idx_color) {
                        cubes[i].idx_color = n;
                        cubes[i].color[0] = r;
                        cubes[i].color[1] = g;
                        cubes[i].color[2] = b;
                    }
                }
            }
        }
        return "";
    }

}
