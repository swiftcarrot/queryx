// Code generated by queryx, DO NOT EDIT.

import { Clause } from "./clause";

export function newDelete(from: string) {
  return new DeleteStatemnet(from);
}

export class DeleteStatemnet {
  private _from: string;
  private _where?: Clause;

  constructor(from: string) {
    this._from = from;
  }

  where(expr: Clause) {
    this._where = expr;
    return this;
  }

  toSQL(): [string, any[]] {
    let sql = "";
    let args: any[] = [];

    sql = `DELETE FROM ${this._from}`;

    if (this._where !== undefined) {
      sql = `${sql} WHERE ${this._where.fragment}`;
      args = args.concat(this._where.args);
    }

    return [sql, args];
  }
}
