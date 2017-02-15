package nl.xs4all.pebbe.cubetalk;

import android.content.ContentValues;
import android.content.Context;
import android.database.Cursor;
import android.database.sqlite.SQLiteDatabase;
import android.database.sqlite.SQLiteOpenHelper;

public class MyDBHandler extends SQLiteOpenHelper {

    private static final int DATABASE_VERSION = 2;
    private static final String DATABASE_NAME = "settings.db";
    public static final String TABLE_SETTINGS = "settings";

    public static final String COLUMN_SETTING = "setting";
    public static final String COLUMN_VALUE = "value";

    public MyDBHandler(Context context) {
        super(context, DATABASE_NAME, null, DATABASE_VERSION);
    }

    @Override
    public void onCreate(SQLiteDatabase db) {
        String CREATE_SETTINGS_TABLE = "CREATE TABLE " +
                TABLE_SETTINGS + "(" +
                COLUMN_SETTING + " TEXT," +
                COLUMN_VALUE + " TEXT)";
        db.execSQL(CREATE_SETTINGS_TABLE);
    }

    @Override
    public void onUpgrade(SQLiteDatabase db, int oldVersion, int newVersion) {
        db.execSQL("DROP TABLE IF EXISTS " + TABLE_SETTINGS);
        onCreate(db);
    }

    public void addSetting(String setting, String value) {
        ContentValues values = new ContentValues();
        values.put(COLUMN_SETTING, setting);
        values.put(COLUMN_VALUE, value);

        SQLiteDatabase db = this.getWritableDatabase();

        db.delete(TABLE_SETTINGS, COLUMN_SETTING + "= \"" + setting + "\"", null);
        db.insert(TABLE_SETTINGS, null, values);
        db.close();
    }

    public String findSetting(String setting) {
        String query = "SELECT " + COLUMN_VALUE + " FROM " + TABLE_SETTINGS + " WHERE " +
                COLUMN_SETTING + " = \"" + setting + "\"";

        SQLiteDatabase db = this.getWritableDatabase();

        Cursor cursor = db.rawQuery(query, null);

        String value = "";

        if (cursor.moveToFirst()) {
            cursor.moveToFirst();
            value = cursor.getString(0);
            cursor.close();
        }
        db.close();
        return value;
    }
}
