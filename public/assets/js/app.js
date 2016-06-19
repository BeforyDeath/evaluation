var Records = (function () {
    function Records() {
        var _this = this;
        this.data = {};
        this.load = function (url) {
            var self = _this;
            $.ajax({
                url: url,
                dataType: 'json',
                async: false,
                success: function (data) {
                    self.data = data;
                }
            });
        };
        this.getSecret = function () {
            var data = false;
            $.ajax({
                url: "api/secret",
                dataType: 'json',
                async: false,
                success: function (_data) {
                    data = _data;
                },
            });
            return data;
        };
    }
    return Records;
}());
var Table = (function () {
    function Table() {
        var _this = this;
        this.CreateTable = function (selector) {
            _this.selector = selector;
            $(selector).append($("<table><thead><tr><tbody>"));
        };
        this.DrawData = function (records) {
            var self = _this;
            var filterMenu = $('#filter-menu');
            if (!records) {
            }
            else {
                $.each(records, function (col, value) {
                    var th = $('<th/>', { title: value.Name });
                    th.attr("data-column", col);
                    if (!value.Params.Sort)
                        th.attr("data-sorter", "false");
                    var title = value.Name;
                    if (value.ShortName)
                        title = value.ShortName;
                    if (!value.Params.Filter)
                        th.attr("data-filter", "false");
                    if (value.Params.Filter || value.Params.Enum) {
                        var li = $('<li>');
                        var a = $("<a>", {
                            href: "#",
                            "data-toggle": "modal",
                            "data-target": "#filter-adding",
                            "data-column": col
                        }).text(title);
                        if (value.Params.Enum)
                            a.attr("data-type", "enum");
                        else
                            a.attr("data-type", "match");
                        a.appendTo(li);
                        filterMenu.append(li);
                        th.attr("data-filter", 'true');
                    }
                    if (value.Params.Priority)
                        th.attr("data-priority", value.Params.Priority);
                    value.Params.Hide = !value.Params.Hide;
                    th.attr("data-columnSelector", value.Params.Hide);
                    th.text(title);
                    $(self.selector + " thead tr").append(th);
                });
                var rows = records[0].Data.length;
                for (var row = 0; row < rows; row++) {
                    var tr = $(self.selector + " tbody").find('[data-row="' + row + '"]');
                    if (tr.length == 0)
                        tr = $('<tr>').attr('data-row', row);
                    for (var col in records) {
                        var item = records[col].Data[row];
                        var value = item.Value;
                        if (item.Link)
                            value = $('<a>', { href: item.Link, target: "_blank" }).text(item.Value);
                        var td = $('<td/>').html(value);
                        tr.append(td);
                    }
                    $(self.selector + " tbody").append(tr);
                }
            }
        };
        this.DrawFilterMenu = function (selectr) {
        };
        this.DrawEnum = function (column, records) {
            var select = $('<select>').addClass('form-control select-enum').attr('multiple', '');
            var result = [];
            for (var i in records[column].Data)
                result[records[column].Data[i].Value] = 1;
            for (var value in result) {
                var option = $('<option/>').text(value);
                select.append(option);
            }
            return select;
        };
    }
    return Table;
}());
var App = (function () {
    function App() {
        this.table = new Table();
        this.records = new Records();
        this.records.load("data/table.json");
    }
    return App;
}());
