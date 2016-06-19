/// <reference path="../typings/jquery/jquery.d.ts" />

class Table {

    selector:string;

    CreateTable = (selector:string) => {
        this.selector = selector;
        $(selector).append($("<table><thead><tr><tbody>"));
    };

    DrawData = (records) => {
        let self = this;
        let filterMenu = $('#filter-menu');
        if (!records) {
        } else {
            $.each(records, function (col, value) {
                let th = $('<th/>', {title: value.Name});

                th.attr("data-column", col);
                if (!value.Params.Sort) th.attr("data-sorter", "false");

                let title = value.Name;
                if (value.ShortName) title = value.ShortName;

                if (!value.Params.Filter) th.attr("data-filter", "false");
                if (value.Params.Filter || value.Params.Enum) {
                    let li = $('<li>');
                    let a = $("<a>", {
                        href: "#",
                        "data-toggle": "modal",
                        "data-target": "#filter-adding",
                        "data-column": col
                    }).text(title);
                    if (value.Params.Enum) a.attr("data-type", "enum");
                    else a.attr("data-type", "match");

                    a.appendTo(li);
                    filterMenu.append(li);
                    th.attr("data-filter", 'true');
                }

                if (value.Params.Priority) th.attr("data-priority", value.Params.Priority);

                value.Params.Hide = !value.Params.Hide;
                th.attr("data-columnSelector", value.Params.Hide);

                th.text(title);
                $(self.selector + " thead tr").append(th);
            });

            let rows = records[0].Data.length;

            for (let row = 0; row < rows; row++) {
                let tr = $(self.selector + " tbody").find('[data-row="' + row + '"]');
                if (tr.length == 0) tr = $('<tr>').attr('data-row', row);

                for (let col in records) {
                    let item = records[col].Data[row];
                    let value = item.Value;
                    if (item.Link) value = $('<a>', {href: item.Link, target: "_blank"}).text(item.Value);
                    let td = $('<td/>').html(value);
                    tr.append(td);
                }

                $(self.selector + " tbody").append(tr);
            }
        }
    };

    DrawFilterMenu = (selectr:string)=> {

    };

    DrawEnum = (column:number, records:any) => {
        let select = $('<select>').addClass('form-control select-enum').attr('multiple', '');
        let result = [];
        for (let i in records[column].Data) result[records[column].Data[i].Value] = 1;
        for (let value in result) {
            let option = $('<option/>').text(value);
            select.append(option);
        }

        return select
    };

}