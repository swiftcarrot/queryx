import { Clause } from "./clause";

export const newTable = (name: string) => {
  return new Table(name);
};

export class Table {
  name: string;

  constructor(name: string) {
    this.name = name;
  }

  newBigIntColumn(name: string) {
    return new BigIntColumn(this, name);
  }

  // TODO: on demand type generation based on schema
  newIntegerColumn(name: string) {
    return new IntegerColumn(this, name);
  }

  newStringColumn(name: string) {
    return new StringColumn(this, name);
  }

  newTextColumn(name: string) {
    return new TextColumn(this, name);
  }

  newDateColumn(name: string) {
    return new DateColumn(this, name);
  }
  newTimeColumn(name: string) {
    return new TimeColumn(this, name);
  }
  newDatetimeColumn(name: string) {
    return new DatetimeColumn(this, name);
  }
  newBooleanColumn(name: string) {
    return new BooleanColumn(this, name);
  }
  newJSONColumn(name: string) {
    return new JSONColumn(this, name);
  }
  newUUIDColumn(name: string) {
    return new UUIDColumn(this, name);
  }
  newFloatColumn(name: string) {
    return new FloatColumn(this, name);
  }
}

class Column {
  table: Table;
  name: string;

  constructor(table: Table, name: string) {
    this.table = table;
    this.name = name;
  }
}

// TODO: on demand generation based on types in schema

class NumberColumn extends Column {
  eq(v: number): Clause {
    return new Clause(`${this.table.name}.${this.name} = ?`, [v]);
  }
}
export class BigIntColumn extends NumberColumn {}
export class IntegerColumn extends NumberColumn {}
export class FloatColumn extends NumberColumn {}

export class BooleanColumn extends Column {}

export class StringColumn extends Column {
  eq(v: string): Clause {
    return new Clause(`${this.table.name}.${this.name} = ?`, [v]);
  }
}

export class TextColumn extends StringColumn {}

export class DateColumn extends Column {}
export class TimeColumn extends Column {}
export class DatetimeColumn extends Column {}
export class UUIDColumn extends Column {}
export class JSONColumn extends Column {}
