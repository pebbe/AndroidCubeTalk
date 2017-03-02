package nl.xs4all.pebbe.cubetalk;

import android.content.DialogInterface;
import android.content.Intent;
import android.os.Bundle;
import android.os.Handler;
import android.os.Message;
import android.support.v7.app.AlertDialog;
import android.support.v7.app.AppCompatActivity;
import android.view.View;
import android.widget.AdapterView;
import android.widget.ArrayAdapter;
import android.widget.EditText;
import android.widget.SeekBar;
import android.widget.Spinner;
import android.widget.TextView;

import java.io.DataInputStream;
import java.io.PrintStream;
import java.net.InetSocketAddress;
import java.net.Socket;

public class StartActivity extends AppCompatActivity {

    private int delay = 2;
    private int enhance = 6;

    static final int VR_REQUEST = 1;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_start);

        String s = getValue(Util.kAddress);
        TextView tv = (TextView) findViewById(R.id.opt_server_address);
        tv.setText(s);

        int i = getIntValue(Util.kPort);
        if (i > 0) {
            tv = (TextView) findViewById(R.id.opt_server_port);
            s = "" + i;
            tv.setText(s);
        }

        s = getValue(Util.kUid);
        tv = (TextView) findViewById(R.id.opt_uid);
        tv.setText(s);

    }

    @Override
    protected void onStop() {
        super.onStop();
        saveServerValues();
    }

    private void saveServerValues() {
        TextView tv = (TextView) findViewById(R.id.opt_server_address);
        String s = tv.getText().toString().replaceAll("\"", "").trim();
        saveValue(Util.kAddress, s);
        tv = (TextView) findViewById(R.id.opt_server_port);
        s = tv.getText().toString().trim();
        saveValue(Util.kPort, s);
        tv = (TextView) findViewById(R.id.opt_uid);
        s = tv.getText().toString().trim();
        saveValue(Util.kUid, s);
    }

    private void saveValue(String key, String value) {
        MyDBHandler handler = new MyDBHandler(this);
        handler.addSetting(key, value);
    }

    private String getValue(String key) {
        MyDBHandler handler = new MyDBHandler(this);
        return handler.findSetting(key);
    }

    private int getIntValue(String key) {
        String s = getValue(key);
        int i = -1;
        if (! s.equals("")) {
            i = Integer.parseInt(s, 10);
        }
        return i;
    }

    Handler runHandler = new Handler() {
        @Override
        public void handleMessage(Message msg) {
            super.handleMessage(msg);
            Bundle bundle = msg.getData();
            if (bundle != null) {
                String e = bundle.getString(Util.sError, "");
                if (!e.equals("")) {
                    alert(e);
                    return;
                }
            }
            runNow();
        }
    };

    public void runNow() {
        Intent i = new Intent(this, MainActivity.class);
        startActivityForResult(i, VR_REQUEST);
    }

    public void run(View view) {
        saveServerValues();

        Runnable runnable = new Runnable() {
            @Override
            public void run() {
                String err = "";

                // TODO: check for wifi
                try {
                    String addr = getValue(Util.kAddress);
                    if (addr.equals("")) {
                        throw new Error("Missing domain or address");
                    }
                    int port = getIntValue(Util.kPort);
                    if (port < 0) {
                        throw new Error("Missing port number");
                    }
                    String uid = getValue(Util.kUid);
                    if (uid.equals("")) {
                        throw new Error("Missing UID");
                    }
                    Socket socket = new Socket();
                    socket.connect(new InetSocketAddress(addr, port), 2000);
                    DataInputStream input = new DataInputStream(socket.getInputStream());
                    PrintStream output = new PrintStream(socket.getOutputStream());
                    output.format("join %s\n", uid);
                    String result = input.readLine().trim();
                    if (!result.equals(".")) {
                        throw new Error("Invalid response from server: " + result);
                    }
                    output.format("quit\n");
                    socket.close();
                } catch (Exception | Error e) {
                    err = e.toString();
                }

                Message msg = Message.obtain();
                Bundle bundle = new Bundle();
                bundle.putString(Util.sError, err);
                msg.setData(bundle);
                runHandler.sendMessage(msg);
            }
        };
        Thread myThread = new Thread(runnable);
        myThread.start();
    }

    public void alert(String err) {
        AlertDialog.Builder builder = new AlertDialog.Builder(this);
        builder.setMessage(err)
                .setTitle(R.string.error)
                .setPositiveButton(R.string.ok, new DialogInterface.OnClickListener() {
                    @Override
                    public void onClick(DialogInterface dialog, int id) {
                    }
                });
        builder.show();
    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        // Check which request we're responding to
        if (requestCode == VR_REQUEST) {
            // Make sure the request was successful
            if (resultCode == RESULT_OK) {
                if (data.hasExtra(Util.sError)) {
                    alert(data.getExtras().getString(Util.sError));
                }
            }
        } else {
            super.onActivityResult(requestCode, resultCode, data);
        }
    }
}
