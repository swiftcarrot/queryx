export const newInsert = () => {
  return new InsertStatement();
};

export class InsertStatement {
  private _into: string;
  private _columns: string[];
  private _values: any[][];
  private _returning: string[];
  private _onConflict: string;

  constructor() {
    // TODO: fix this
    this._into = "";
    this._columns = [];
    this._values = [];
    this._returning = [];
    this._onConflict = "";
  }

  into(into: string) {
    this._into = into;
    return this;
  }

  columns(...columns: string[]) {
    this._columns = columns;
    return this;
  }

  values(...values: any[]) {
    this._values.push(values);
    return this;
  }

  returning(...returning: string[]) {
    this._returning = returning;
    return this;
  }

  onConflict(onConflict: string) {
    this._onConflict = onConflict;
    return this;
  }

  toSQL(): [string, any[]] {
    let sql: string = `INSERT INTO ${this._into}`;

    if (this._columns.length > 0) {
      sql = `${sql} (${this._columns.join(", ")})`;
    }

    const values: string[] = [];
    for (const v of this._values) {
      const ss: string[] = [];
      for (let i = 0; i < v.length; i++) {
        ss.push("?");
      }
      values.push(`(${ss.join(", ")})`);
    }
    sql = `${sql} VALUES ${values.join(", ")}`;

    if (this._returning.length > 0) {
      sql = `${sql} RETURNING ${this._returning.join(", ")}`;
    }

    if (this._onConflict !== "") {
      sql = `${sql} ${this._onConflict}`;
    }

    const args: any[] = [];
    for (const v of this._values) {
      args.push(...v);
    }

    return [sql, args];
  }
}
