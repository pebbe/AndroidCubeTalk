package nl.xs4all.pebbe.cubetalk;

import android.content.Intent;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.opengl.GLES20;
import android.opengl.GLUtils;
import android.opengl.Matrix;
import android.os.Bundle;

import com.google.vr.sdk.base.Eye;
import com.google.vr.sdk.base.GvrActivity;
import com.google.vr.sdk.base.GvrView;
import com.google.vr.sdk.base.HeadTransform;
import com.google.vr.sdk.base.Viewport;

import java.io.DataInputStream;
import java.io.PrintStream;
import java.net.Socket;
import java.util.Locale;

import javax.microedition.khronos.egl.EGLConfig;

public class MainActivity extends GvrActivity implements GvrView.StereoRenderer {

    private static final int NR_OF_CONNECTIONS = 16;

    private static final int MAX_CUBES = 100;

    private int currentConnection = 0;

    private int syncNrOfCubes = 0;

    private Socket[] sockets;
    private DataInputStream[] inputs;
    private PrintStream[] outputs;

    private String[] ids;
    private CubeData[] cubes;
    private float syncSelfZ = 4;
    private float selfZ = -4;
    private long selfIdx = 0;
    private long syncInfoIdx = 0;
    private String syncInfoID = "";
    private String infoID = "";
    private int infoChoice = 0;
    private String syncInfoChoice1 = "";
    private String infoChoice1 = "";
    private String syncInfoChoice2 = "";
    private String infoChoice2 = "";
    private boolean syncHasInfo = false;
    private boolean syncHasChoice = false;
    private boolean hasInfo = false;
    private boolean hasChoice = false;
    private boolean syncReplyChoice = false;
    private boolean syncMark = false;
    private String syncReplyChoiceID = "";
    private String syncReplyChoiceText = "";
    private String[] syncInfoLines;
    private Info info;
    private long syncCubeSizeIdx = 0;
    private boolean syncSetCubeSize = false;
    private float[] syncCubeSize;
    private float[] cubeSize;
    private boolean syncErr = false;
    private String syncErrStr = "";
    final private Object settingsLock = new Object();

    private boolean[] runnings;
    final private Object runningLock = new Object();

    private int nrOfCubes;

    private Kubus kubus;
    private Wereld wereld;
    private int[] texturenames;

    protected float[] modelCube;
    protected float[] modelWorld;
    protected float[] modelInfo;
    private float[] camera;
    private float[] view;
    private float[] modelViewProjection;
    private float[] modelView;
    private float[] forward;
    private float infoAngleH;
    private float infoAngleV;

    private GvrView gvrView;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);

        initializeGvrView();

        // Initialize other objects here.
        modelCube = new float[16];
        modelWorld = new float[16];
        modelInfo = new float[16];
        camera = new float[16];
        view = new float[16];
        modelViewProjection = new float[16];
        modelView = new float[16];
        forward = new float[3];

        syncCubeSize = new float[3];
        cubeSize = new float[3];
        cubeSize[0] = 1;
        cubeSize[1] = 1;
        cubeSize[2] = 1;

        MyDBHandler handler = new MyDBHandler(this);
        // external server
        String address = handler.findSetting(Util.kAddress);
        String p = handler.findSetting(Util.kPort);
        int port = 0;
        if (!p.equals("")) {
            port = Integer.parseInt(p, 10);
        }
        String id = handler.findSetting(Util.kUid);

        ids = new String[MAX_CUBES];
        cubes = new CubeData[MAX_CUBES];

        sockets = new Socket[NR_OF_CONNECTIONS];
        inputs = new DataInputStream[NR_OF_CONNECTIONS];
        outputs = new PrintStream[NR_OF_CONNECTIONS];
        runnings = new boolean[NR_OF_CONNECTIONS];
        for (int i = 0; i < NR_OF_CONNECTIONS; i++) {
            runnings[i] = true;
        }

        final String uid = id;
        final String addr = address;
        final int pnum = port;

        Runnable runnable = new Runnable() {
            @Override
            public void run() {
                for (int i = 0; i < NR_OF_CONNECTIONS; i++) {
                    try {
                        sockets[i] = new Socket(addr, pnum);
                        sockets[i].setSoTimeout(10000);
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
    public void onSurfaceCreated(EGLConfig config) {
        Matrix.setLookAtM(camera, 0,
                0.0f, 0.0f, 0.01f,  // 0.01f
                0.0f, 0.0f, 0.0f,
                0.0f, 1.0f, 0.0f);

        texturenames = new int[Util.NR_OF_TEXTURES];
        GLES20.glGenTextures(Util.NR_OF_TEXTURES, texturenames, 0);

        setBitMap(R.raw.head5a, Util.TEXTURE_HEAD0);
        setBitMap(R.raw.head5b, Util.TEXTURE_HEAD1);
        setBitMap(R.raw.head5c, Util.TEXTURE_HEAD2);
        setBitMap(R.raw.face1a, Util.TEXTURE_FACE0);
        setBitMap(R.raw.face1b, Util.TEXTURE_FACE1);
        setBitMap(R.raw.face1c, Util.TEXTURE_FACE2);

        kubus = new Kubus();
        wereld = new Wereld(this, texturenames[Util.TEXTURE_WORLD]);
    }

    @Override
    public void onNewFrame(HeadTransform headTransform) {

        synchronized (settingsLock) {
            selfZ = syncSelfZ;
            if (syncSetCubeSize) {
                syncSetCubeSize = false;
                cubeSize[0] = syncCubeSize[0];
                cubeSize[1] = syncCubeSize[1];
                cubeSize[2] = syncCubeSize[2];
            }
        }

        Matrix.setIdentityM(modelWorld, 0);
        Matrix.setIdentityM(modelInfo, 0);

        Matrix.translateM(modelWorld, 0, 0, 0, -selfZ);

        headTransform.getForwardVector(forward, 0);
        float[] euler = new float[3];
        headTransform.getEulerAngles(euler, 0);
        int retval = doForward(forward, euler[2] / (float) Math.PI * 180.0f);
        if (retval == Util.stERROR) {
            Intent data = new Intent();
            String e;
            synchronized (settingsLock) {
                e = syncErrStr;
            }
            data.putExtra(Util.sError, e);
            setResult(RESULT_OK, data);
            finish();
            return;
        }
        synchronized (settingsLock) {
            nrOfCubes = syncNrOfCubes;
            if (syncHasInfo) {
                syncHasInfo = false;
                hasInfo = true;
                hasChoice = syncHasChoice;
                syncHasChoice = false;
                infoChoice1 = syncInfoChoice1;
                infoChoice2 = syncInfoChoice2;
                infoID = syncInfoID;
                infoAngleH = (float) (-Math.atan2(forward[0], -forward[2]) / Math.PI * 180.0);
                infoAngleV = (float) (Math.atan2(forward[1], Math.sqrt(forward[0] * forward[0] + forward[2] * forward[2])) / Math.PI * 180.0);
                info = new Info(this, texturenames[Util.TEXTURE_INFO0], texturenames[Util.TEXTURE_INFO1], hasChoice, infoChoice1, infoChoice2, syncInfoLines);
            }
        }
        if (hasInfo) {
            Matrix.setIdentityM(modelInfo, 0);
            Matrix.rotateM(modelInfo, 0, infoAngleH, 0, 1, 0);
            Matrix.rotateM(modelInfo, 0, infoAngleV, 1, 0, 0);
            Matrix.translateM(modelInfo, 0, 0, 0, -selfZ);
            //Matrix.rotateM(modelInfo, 0, infoAngle, 0, 1, 0);

            infoChoice = 0;
            if (hasChoice) {
                float roth = (float) (-Math.atan2(forward[0], -forward[2]) / Math.PI * 180.0) - infoAngleH;
                if (roth < -180) {
                    roth += 360;
                } else if (roth > 180) {
                    roth -= 360;
                }
                if (roth < 0) {
                    infoChoice = 1;
                }
            }

        }
    }

    @Override
    public void onDrawEye(Eye eye) {
        GLES20.glClear(GLES20.GL_COLOR_BUFFER_BIT | GLES20.GL_DEPTH_BUFFER_BIT);

        // Apply the eye transformation to the camera.
        Matrix.multiplyMM(view, 0, eye.getEyeView(), 0, camera, 0);

        // Build the ModelView and ModelViewProjection matrices
        // for calculating cube position and light.
        float[] perspective = eye.getPerspective(0.1f, 200.0f);

        Matrix.multiplyMM(modelView, 0, view, 0, modelWorld, 0);
        Matrix.multiplyMM(modelViewProjection, 0, perspective, 0, modelView, 0);

        GLES20.glDisable(GLES20.GL_CULL_FACE);
        wereld.draw(modelViewProjection);

        GLES20.glEnable(GLES20.GL_CULL_FACE);
        GLES20.glCullFace(GLES20.GL_BACK);
        GLES20.glEnable(GLES20.GL_DEPTH_TEST);
        GLES20.glDepthFunc(GLES20.GL_LEQUAL);


        int n;

        for (int i = 0; i < nrOfCubes; i++) {
            float red;
            float green;
            float blue;
            int texturehead;
            int textureface;
            synchronized (settingsLock) {
                if (!cubes[i].visible) {
                    continue;
                }

                Matrix.setIdentityM(modelCube, 0);

                Matrix.translateM(modelCube, 0, 0, 0, -selfZ);
                Matrix.translateM(modelCube, 0, cubes[i].pos[0], cubes[i].pos[1], cubes[i].pos[2]);

                float x = cubes[i].lookat[0];
                float y = cubes[i].lookat[1];
                float z = cubes[i].lookat[2];

                float roth = (float) Math.atan2(x, z);
                Matrix.rotateM(modelCube, 0, roth / (float) Math.PI * 180.0f, 0, 1, 0);
                float rotv = (float) Math.atan2(y, Math.sqrt(x * x + z * z));
                Matrix.rotateM(modelCube, 0, -rotv / (float) Math.PI * 180.0f, 1, 0, 0);

                Matrix.rotateM(modelCube, 0, cubes[i].roll, 0, 0, -1);

                Matrix.multiplyMM(modelView, 0, view, 0, modelCube, 0);
                Matrix.multiplyMM(modelViewProjection, 0, perspective, 0, modelView, 0);

                red = cubes[i].color[0];
                green = cubes[i].color[1];
                blue = cubes[i].color[2];

                texturehead = texturenames[cubes[i].texturehead];
                textureface = texturenames[cubes[i].textureface];
            }
            kubus.draw(modelViewProjection, red, green, blue, cubeSize, texturehead, textureface);
        }

        if (hasInfo) {
            GLES20.glEnable(GLES20.GL_BLEND);
            Matrix.multiplyMM(modelView, 0, view, 0, modelInfo, 0);
            Matrix.multiplyMM(modelViewProjection, 0, perspective, 0, modelView, 0);
            info.draw(modelViewProjection, infoChoice);
            GLES20.glDisable(GLES20.GL_BLEND);
        }

    }

    @Override
    public void onRendererShutdown() {

    }

    @Override
    public void onSurfaceChanged(int i, int i1) {

    }

    @Override
    public void onFinishFrame(Viewport viewport) {

    }

    public void initializeGvrView() {
        setContentView(R.layout.common_ui);

        gvrView = (GvrView) findViewById(R.id.gvr_view);
        gvrView.setEGLConfigChooser(8, 8, 8, 8, 16, 8);

        gvrView.setRenderer(this);
        //gvrView.setTransitionViewEnabled(true);

        // Enable Cardboard-trigger feedback with Daydream headsets. This is a simple way of supporting
        // Daydream controller input for basic interactions using the existing Cardboard trigger API.
        gvrView.enableCardboardTriggerEmulation();

        /*
        if (gvrView.setAsyncReprojectionEnabled(true)) {
            // Async reprojection decouples the app framerate from the display framerate,
            // allowing immersive interaction even at the throttled clockrates set by
            // sustained performance mode.
            AndroidCompat.setSustainedPerformanceMode(this, true);
        }
        */

        setGvrView(gvrView);
    }

    @Override
    public void onCardboardTrigger() {
        if (hasInfo) {
            if (hasChoice) {
                synchronized (settingsLock) {
                    syncReplyChoice = true;
                    syncReplyChoiceID = infoID;
                    syncReplyChoiceText = infoChoice == 0 ? infoChoice1 : infoChoice2;
                }
                hasChoice = false;
            }
            hasInfo = false;
            info = null;
        } else {
            synchronized (settingsLock) {
                syncMark = true;
            }
        }
    }

    private int doForward(float[] in, float roll) {
        final float xi = in[0];
        final float yi = in[1];
        final float zi = in[2];
        final float ri = roll;
        final int index = currentConnection;
        currentConnection = (currentConnection + 1) % NR_OF_CONNECTIONS;
        Runnable runnable = new Runnable() {
            @Override
            public void run() {
                synchronized (runningLock) {
                    if (runnings[index]) {
                        return;
                    }
                    runnings[index] = true;
                }
                boolean replyChoice = false;
                String replyID = "";
                String replyText = "";
                String replyMark = "";
                synchronized (settingsLock) {
                    replyChoice = syncReplyChoice;
                    if (replyChoice) {
                        replyID = syncReplyChoiceID;
                        replyText = syncReplyChoiceText;
                        syncReplyChoice = false;
                    } else if (syncMark) {
                        syncMark = false;
                        replyMark = "mark";
                    }
                }
                if (replyChoice) {
                    outputs[index].format(Locale.US, "info %s %s\n", replyID, replyText);
                } else {
                    outputs[index].format(Locale.US, "lookat %f %f %f %f %s\n", xi, yi, zi, ri, replyMark);
                }

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
                            syncErr = true;
                            syncErrStr = e.toString();
                        }
                        return;
                    }

                    if (response == null) {
                        syncErr = true;
                        syncErrStr = "No response from remote server";
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
                        } else if (parts[0].equals("head")) {
                            e = setHead(parts);
                        } else if (parts[0].equals("face")) {
                            e = setFace(parts);
                        } else if (parts[0].equals("info")) {
                            e = setInfo(parts, index);
                        } else if (parts[0].equals("cubesize")) {
                            e = setCubesize(parts, index);
                        } else if (parts[0].equals("recenter")) {
                            gvrView.recenterHeadTracker();
                        }
                    }
                    if (!e.equals("")) {
                        syncErr = true;
                        syncErrStr = e;
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
            if (syncErr) {
                retval = Util.stERROR;
                syncErr = false;
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
                if (n >= selfIdx) {
                    selfIdx = n;
                    syncSelfZ = z;
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
                for (i = 0; i < syncNrOfCubes; i++) {
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
                    if (syncNrOfCubes < MAX_CUBES) {
                        cubes[syncNrOfCubes] = new CubeData(Util.TEXTURE_HEAD0, Util.TEXTURE_FACE0);
                        cubes[syncNrOfCubes].visible = true;
                        ids[i] = parts[1];
                        syncNrOfCubes++;
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
                for (i = 0; i < syncNrOfCubes; i++) {
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
                for (i = 0; i < syncNrOfCubes; i++) {
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

    // lookat {id} {n3} {x} {y} {z} {r}
    private String setLookat(String[] parts) {
        if (parts.length == 7) {
            synchronized (settingsLock) {
                long n = 0;
                float x;
                float y;
                float z;
                float r;
                try {
                    n = Integer.parseInt(parts[2]);
                    x = Float.parseFloat(parts[3]);
                    y = Float.parseFloat(parts[4]);
                    z = Float.parseFloat(parts[5]);
                    r = Float.parseFloat(parts[6]);
                } catch (Exception e) {
                    return e.toString();
                }
                boolean found = false;
                int i;
                for (i = 0; i < syncNrOfCubes; i++) {
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
                        cubes[i].roll = r;
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
                for (i = 0; i < syncNrOfCubes; i++) {
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

    // head {id} {n7} {head}
    private String setHead(String[] parts) {
        if (parts.length == 4) {
            synchronized (settingsLock) {
                long n = 0;
                int h;
                try {
                    n = Integer.parseInt(parts[2]);
                    h = Integer.parseInt(parts[3]);
                } catch (Exception e) {
                    return e.toString();
                }
                boolean found = false;
                int i;
                for (i = 0; i < syncNrOfCubes; i++) {
                    if (ids[i].equals(parts[1])) {
                        found = true;
                        break;
                    }
                }
                if (found) {
                    if (n >= cubes[i].idx_head) {
                        cubes[i].idx_head = n;
                        cubes[i].texturehead = getHeadTexture(h);
                    }
                }
            }
        }
        return "";
    }

    // face {id} {n7} {face}
    private String setFace(String[] parts) {
        if (parts.length == 4) {
            synchronized (settingsLock) {
                long n = 0;
                int h;
                try {
                    n = Integer.parseInt(parts[2]);
                    h = Integer.parseInt(parts[3]);
                } catch (Exception e) {
                    return e.toString();
                }
                boolean found = false;
                int i;
                for (i = 0; i < syncNrOfCubes; i++) {
                    if (ids[i].equals(parts[1])) {
                        found = true;
                        break;
                    }
                }
                if (found) {
                    if (n >= cubes[i].idx_face) {
                        cubes[i].idx_face = n;
                        cubes[i].textureface = getFaceTexture(h);
                    }
                }
            }
        }
        return "";
    }


    // info {n5} {nr of lines}
    // info {n5} {nr of lines} {responce ID} {choice 1} {choice 2}
    private String setInfo(String[] parts, int index) {
        if (parts.length == 3 || parts.length == 6) {
            long n = 0;
            int nr_of_lines = 0;
            try {
                n = Integer.parseInt(parts[1]);
                nr_of_lines = Integer.parseInt(parts[2]);
            } catch (Exception e) {
                return e.toString();
            }
            String[] lines = new String[nr_of_lines];
            for (int i = 0; i < nr_of_lines; i++) {
                try {
                    lines[i] = inputs[index].readLine();
                } catch (Exception e) {
                    return e.toString();
                }
            }
            if (n >= syncInfoIdx) {
                synchronized (settingsLock) {
                    syncInfoIdx = n;
                    syncHasInfo = true;
                    syncInfoID = parts[2];
                    syncInfoLines = new String[nr_of_lines];
                    for (int i = 0; i < nr_of_lines; i++) {
                        syncInfoLines[i] = lines[i];
                    }

                    if (parts.length == 3) {
                        syncHasChoice = false;
                    } else {
                        syncHasChoice = true;
                        syncInfoID = parts[3];
                        syncInfoChoice1 = parts[4];
                        syncInfoChoice2 = parts[5];
                    }
                }
            }
        }
        return "";
    }

    // cubesize {n6} {w} {h} {d}
    private String setCubesize(String[] parts, int index) {
        if (parts.length == 5) {
            long n = 0;
            float w = 1;
            float h = 1;
            float d = 1;
            try {
                n = Integer.parseInt(parts[1]);
                w = Float.parseFloat(parts[2]);
                h = Float.parseFloat(parts[3]);
                d = Float.parseFloat(parts[4]);
            } catch (Exception e) {
                return e.toString();
            }
            synchronized (settingsLock) {
                if (n >= syncCubeSizeIdx) {
                    syncCubeSizeIdx = n;
                    syncCubeSize[0] = w;
                    syncCubeSize[1] = h;
                    syncCubeSize[2] = d;
                    syncSetCubeSize = true;
                }
            }
        }
        return "";
    }

    private int getHeadTexture(int h) {
        if (h == 1) return Util.TEXTURE_HEAD1;
        if (h == 2) return Util.TEXTURE_HEAD2;
        return Util.TEXTURE_HEAD0;
    }

    private int getFaceTexture(int h) {
        if (h == 1) return Util.TEXTURE_FACE1;
        if (h == 2) return Util.TEXTURE_FACE2;
        return Util.TEXTURE_FACE0;
    }

    private void setBitMap(int bitmap, int texture) {

        // Temporary create a bitmap
        Bitmap bmp = BitmapFactory.decodeResource(getResources(), bitmap);

        // Bind texture to texturename
        GLES20.glActiveTexture(GLES20.GL_TEXTURE0);
        Util.checkGlError("glActiveTexture");
        GLES20.glBindTexture(GLES20.GL_TEXTURE_2D, texturenames[texture]);
        Util.checkGlError("glBindTexture");

        // Load the bitmap into the bound texture.
        GLUtils.texImage2D(GLES20.GL_TEXTURE_2D, 0, bmp, 0);
        Util.checkGlError("texImage2D");

        // We are done using the bitmap so we should recycle it.
        bmp.recycle();

    }
}