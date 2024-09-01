// Code generated by queryx, DO NOT EDIT.

import Database from "better-sqlite3";
import { format, parse } from "date-fns";
import { Config } from "./config";

export class Adapter {
  public config: Config;
  public db: Database;

  constructor(config: Config) {
    this.config = config;
  }

  connect() {
    this.db = new Database(this.config.url);
  }

  newConnection() {
    return new Database(this.config.url);
  }

  release() {}

  query<R>(query: string, ...args: any[]): R[] {
    let [query1, args1] = rebind(query, args);
    let stmt = this.db.prepare(query1);
    return stmt.all(...args1) as R[];
  }

  queryOne<R>(query: string, ...args: any[]): R {
    let [query1, args1] = rebind(query, args);
    let stmt = this.db.prepare(query1);
    return stmt.get(...args1) as R;
  }

  async exec(query: string, ...args: any[]) {
    let [query1, args1] = rebind(query, args);
    let stmt = this.db.prepare(query1);
    let res = stmt.run(...args1);
    return res.changes;
  }

  async beginTx() {
    await this.exec("BEGIN");
  }

  async commit() {
    await this.exec("COMMIT");
  }

  async rollback() {
    await this.exec("ROLLBACK");
  }
}

export function rebind<T extends any[] = any[]>(query: string, args?: T) {
  let str = "";
  let i = 0;
  let j = 1;
  let k = 0;
  let args1: any[] = [];

  while (i !== -1) {
    i = query.indexOf("?");
    str += query.substring(0, i);

    if (args.length > k) {
      const arg = args[k];
      if (Array.isArray(arg)) {
        args1 = args1.concat(arg);
        str += arg.map((_, i) => "?").join(", ");
        j += arg.length;
      } else {
        args1.push(arg);
        str += "?";
        j++;
      }
      k++;
    }

    query = query.substring(i + 1);
  }

  return [str + query, args1];
}

// convert into sqlite type
export function adapterValue(type: string, value: any) {
  if (typeof value === "object") {
    switch (type) {
      case "time":
        return format(value, "HH:mm:ss");
      case "date":
        return format(value, "yyyy-MM-dd");
      case "datetime":
        return format(value, "yyyy-MM-dd HH:mm:ss");
      case "json":
      case "jsonb":
        return JSON.stringify(value);
      default:
        return value;
    }
  }

  switch (type) {
    case "boolean":
      return value ? 1 : 0;
  }

  return value;
}

// convert into queryx type
export function adapterScan(type: string, value: any) {
  switch (type) {
    case "time":
      return parse(value, "HH:mm:ss", new Date());
    case "date":
      return parse(value, "yyyy-MM-dd", new Date());
    case "datetime":
      return parse(value, "yyyy-MM-dd HH:mm:ss", new Date());
    case "json":
    case "jsonb":
      return JSON.parse(value);
    case "boolean":
      return value === 1;
    default:
      return value;
  }
}
