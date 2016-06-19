/// <reference path="../typings/jquery/jquery.d.ts" />
/// <reference path="Records.ts" />
/// <reference path="Table.ts" />
class App {
    records:Records;
    table:Table;

    constructor() {
        this.table = new Table();
        this.records = new Records();
        this.records.load("data/table.json");
    }
}
