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

public class StartActivity extends AppCompatActivity implements AdapterView.OnItemSelectedListener {

    private int delay = 2;
    private int enhance = 6;

    static final int VR_REQUEST = 1;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_start);

        int i = getIntValue(Util.kMode);
        if (i < 0) {
            i = 0;
        }
        showHideDelayed(i == 0);
        showHideServer(i == 1);

        Spinner optMode = (Spinner) findViewById(R.id.opt_mode);
        ArrayAdapter<CharSequence> adapter = ArrayAdapter.createFromResource(this,
                R.array.mode_array, android.R.layout.simple_spinner_item);
        // Specify the layout to use when the list of choices appears
        adapter.setDropDownViewResource(android.R.layout.simple_spinner_dropdown_item);
        // Apply the adapter to the spinner
        optMode.setAdapter(adapter);
        optMode.setSelection(i);
        optMode.setOnItemSelectedListener(this);

        delay = getIntValue(Util.kDelay);
        if (delay < 0) {
            delay = 2;
        }
        setDelay(delay);

        SeekBar optDelay = (SeekBar) findViewById(R.id.opt_delay);
        optDelay.setProgress(delay);
        optDelay.setOnSeekBarChangeListener(new SeekBar.OnSeekBarChangeListener() {
            @Override
            public void onProgressChanged(SeekBar seekBar, int i, boolean b) {
                setDelay(i);
            }

            @Override
            public void onStartTrackingTouch(SeekBar seekBar) {

            }

            @Override
            public void onStopTrackingTouch(SeekBar seekBar) {
                saveValue(Util.kDelay, ""+delay);
            }
        });

        enhance = getIntValue(Util.kEnhance);
        if (enhance < 0) {
            enhance = 6;
        }
        setEnhance(enhance);

        SeekBar optEnhance = (SeekBar) findViewById(R.id.opt_enhance);
        optEnhance.setProgress(enhance);
        optEnhance.setOnSeekBarChangeListener(new SeekBar.OnSeekBarChangeListener() {
            @Override
            public void onProgressChanged(SeekBar seekBar, int i, boolean b) {
                setEnhance(i);
            }

            @Override
            public void onStartTrackingTouch(SeekBar seekBar) {

            }

            @Override
            public void onStopTrackingTouch(SeekBar seekBar) {
                saveValue(Util.kEnhance, ""+enhance);
            }
        });


        String s = getValue(Util.kAddress);
        TextView tv = (TextView) findViewById(R.id.opt_server_address);
        tv.setText(s);

        i = getIntValue(Util.kPort);
        if (i > 0) {
            tv = (TextView) findViewById(R.id.opt_server_port);
            s = "" + i;
            tv.setText(s);
        }

    }

    @Override
    protected void onStop() {
        super.onStop();
        saveServerValues();
    }

    @Override
    public void onItemSelected(AdapterView<?> parent, View view, int position, long id) {
        saveValue(Util.kMode, ""+position);
        showHideDelayed(position == 0);
        showHideServer(position == 1);
    }

    @Override
    public void onNothingSelected(AdapterView<?> parent) {

    }

    private void setDelay(int i) {
        delay = i;
        TextView tv = (TextView) findViewById(R.id.val_delay);
        String s = String.format("%.1f", 0.5 * (float)i);
        tv.setText(s);
    }

    private void setEnhance(int i) {
        enhance = i;
        TextView tv = (TextView) findViewById(R.id.val_enhance);
        String s = String.format("%.1f", 0.5 * (float)(i-4));
        tv.setText(s);
    }

    private void saveServerValues() {
        TextView tv = (TextView) findViewById(R.id.opt_server_address);
        String s = tv.getText().toString().replaceAll("\"", "").trim();
        saveValue(Util.kAddress, s);
        tv = (TextView) findViewById(R.id.opt_server_port);
        s = tv.getText().toString().trim();
        saveValue(Util.kPort, s);
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

    private void showHideDelayed(boolean show) {
        int v = show ? View.VISIBLE : View.INVISIBLE;
        TextView tv = (TextView) findViewById(R.id.lbl_delay);
        tv.setVisibility(v);
        tv = (TextView) findViewById(R.id.val_delay);
        tv.setVisibility(v);
        SeekBar sb = (SeekBar) findViewById(R.id.opt_delay);
        sb.setVisibility(v);
        tv = (TextView) findViewById(R.id.lbl_enhance);
        tv.setVisibility(v);
        tv = (TextView) findViewById(R.id.val_enhance);
        tv.setVisibility(v);
        sb = (SeekBar) findViewById(R.id.opt_enhance);
        sb.setVisibility(v);
   }

    private void showHideServer(boolean show) {
        int v = show ? View.VISIBLE : View.INVISIBLE;
        TextView tv = (TextView) findViewById(R.id.lbl_server_address);
        tv.setVisibility(v);
        tv = (TextView) findViewById(R.id.lbl_server_port);
        tv.setVisibility(v);
        EditText et = (EditText) findViewById(R.id.opt_server_address);
        et.setVisibility(v);
        et = (EditText) findViewById(R.id.opt_server_port);
        et.setVisibility(v);
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

                switch (getValue(Util.kMode)) {
                    case Util.vModeDelay:
                        break;
                    case Util.vModeExtern:
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
                                uid = "" + System.currentTimeMillis();
                                saveValue(Util.kUid, uid);
                            }
                            Socket socket = new Socket();
                            socket.connect(new InetSocketAddress(addr, port), 2000);
                            DataInputStream input = new DataInputStream(socket.getInputStream());
                            PrintStream output = new PrintStream(socket.getOutputStream());
                            output.format("VRC1.0 %s\n", uid);
                            String result = input.readLine().trim();
                            if (!result.equals("VRC1.0.OK")) {
                                throw new Error("Invalid response from server: " + result);
                            }
                            output.format("quit\n");
                            socket.close();
                        } catch (Exception | Error e) {
                            err = e.toString();
                        }
                        break;
                    default:
                        err = "Invalid mode";
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
