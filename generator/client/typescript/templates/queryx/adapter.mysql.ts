// Code generated by queryx, DO NOT EDIT.

import mysql, { ResultSetHeader, RowDataPacket } from "mysql2/promise";
import { parse } from "date-fns";
import { Config } from "./config";

export class Adapter {
  public config: Config;
  public db: mysql.Pool;

  constructor(config: Config) {
    this.config = config;
  }

  connect() {
    const pool = mysql.createPool({
      uri: this.config.url,
    });
    this.db = pool;
  }

  release() {
    this.db.release();
  }

  async query<R>(query: string, ...args: any[]) {
    let [rows] = await this.db.query<R & RowDataPacket[]>(query, args);
    return rows;
  }

  async queryOne<R>(query: string, ...args: any[]) {
    let [rows] = await this.db.query<R & RowDataPacket[]>(query, args);
    return rows[0] || null;
  }

  async exec(query: string, ...args: any[]) {
    let [res] = await this.db.execute<ResultSetHeader>(query, args);
    return res.affectedRows;
  }

  async _exec(query: string, ...args: any[]) {
    let [res] = await this.db.execute<ResultSetHeader>(query, args);
    return res;
  }

  async beginTx() {
    await this.db.query("START TRANSACTION");
  }

  async commit() {
    await this.db.query("COMMIT");
  }

  async rollback() {
    await this.db.query("ROLLBACK");
  }
}

export function adapterValue(type: string, value: any) {
  switch (type) {
    case "time":
      return parse(value, "HH:mm:ss", new Date());
    default:
      return value;
  }
}

export function adapterScan(type: string, value: any) {
  switch (type) {
    case "boolean":
      return value === 1;
    default:
      return value;
  }
}
