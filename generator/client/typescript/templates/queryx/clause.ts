// Code generated by queryx, DO NOT EDIT.

export class Clause {
  fragment: string;
  args: any[];

  constructor(fragment: string, args: any[]) {
    this.fragment = fragment;
    this.args = args;
  }

  and(...clauses: Clause[]) {
    if (clauses.length === 0) {
      return this;
    }

    let fragment = this.fragment;
    let args = this.args;

    for (let clause of clauses) {
      fragment = `(${this.fragment}) AND (${clause.fragment})`;
      args = this.args.concat(clause.args);
    }
    return new Clause(fragment, args);
  }

  or(...clauses: Clause[]) {
    if (clauses.length === 0) {
      return this;
    }

    let fragment = this.fragment;
    let args = this.args;

    for (let clause of clauses) {
      fragment = `(${this.fragment}) OR (${clause.fragment})`;
      args = this.args.concat(clause.args);
    }
    return new Clause(fragment, args);
  }
}
