var App = new App();
App.table.CreateTable("#table");
App.table.DrawData(App.records.data);

//*/
// todo options tablesorter
$("#table table").tablesorter({
    theme: 'default',

    widthFixed: true,
    widgets: ["zebra", "filter", 'columnSelector', 'columns'],

    //headers: {5: {sorter: false, filter: true}},

    widgetOptions: {
        filter_columnFilters: false,
        //filter_hideFilters: true,

        filter_saveFilters: false,
        filter_resetOnEsc: false,
        filter_reset: '.reset',

        filter_searchFiltered: true, // if true, the search is performed on already filtered rows, with some exceptions

        filter_hideEmpty: true,
        filter_cssFilter: '',
        filter_childRows: false,
        filter_ignoreCase: true,
        filter_searchDelay: 300,
        filter_startsWith: false,

        columnSelector_saveColumns: false,
        columnSelector_mediaquery: true,
        columnSelector_mediaqueryHidden: true,
        columnSelector_breakpoints : [ '20em', '30em', '40em', '50em', '60em', '70em' ],

    }
}).bind("sortEnd", function (e, filter) {
    //console.log(filter);
    //$.tablesorter.getColumnText($("table"), 4, tee, 0);
});
//*/

// todo menu column selector
$.tablesorter.columnSelector.attachTo($('table'), '#popover-target');
$('#popover').popover({
    placement: 'right',
    html: true,
    content: $('#popover-target')
});

// todo show/hidden modal filter adding
$('#filter-adding')
    .on('show.bs.modal', function (e) {
        var column = $(e.relatedTarget).data('column');
        var type = $(e.relatedTarget).data('type');

        $(this).data('column', column);
        $(this).data('type', type);
        $(this).data('name', $(e.relatedTarget).text());


        var body = $('.modal-body-hidden#filter-' + type).children().clone();
        var filterCondition = body.find('.filter-condition');
        if (type == "enum") {
            filterCondition.html(App.table.DrawEnum(column, App.records.data));
        } else filterCondition.html("");


        body.find('ul li a').on('click', AddingCondition);
        body.appendTo($(this).find('.modal-body').html(""));

        $(this).find('.modal-title').text('Добавить фильтр: ' + $(e.relatedTarget).text());
    })
    .on('hidden.bs.modal', function (e) {
        var column = $(this).data('column') + 0;
        var type = $(this).data('type');

        var conditionList = false;
        var filter = "";

        if (type == "enum") {
            conditionList = $(this).find('select').val();
            if (conditionList) {
                for (var value in conditionList) {
                    var title = '=' + conditionList[value];
                    if (filter != "") filter += " or " + title;
                    else filter = title;
                }
            }
        } else if (type == "match") {
            conditionList = $(this).find(".filter-condition .input-group");
            if (conditionList.length > 0) {
                conditionList.each(function (i, elem) {
                    var name = $(elem).data('filter-name');
                    var value = $(elem).find('input').val();
                    if (value != '') {
                        if (name == "more") value = ">" + value;
                        if (name == "less") value = "<" + value;
                        if (name == "equally") value = "=" + value;
                        if (name == "range") value += " - " + $(elem).find('input').next('input').val();
                        if (filter != "") filter += " and " + value;
                        else filter = value;
                    }
                });
            }
        }

        //console.log(column, filter);

        if (filter) {

            // todo обавление фильтра в таблице
            var filters = $.tablesorter.getFilters($("#table table"));
            if (filters[column] != "") filters[column] += " or " + filter;
            else filters[column] = filter;

            $.tablesorter.setFilters($("#table table"), filters, true);


            // todo создание плашки
            var visfilters = $('#vis-filters');
            visfilters.find('[data-column="' + column + '"]').remove();

            var visfilter = $('#vis-filter').find('.vis-filter').clone();

            visfilter.find(".delete").on('click', function () {
                var column = $(this).closest('.vis-filter').data('column');
                var filters = $.tablesorter.getFilters($('table'));
                filters[column] = "";
                $.tablesorter.setFilters($('table'), filters, true);
                $(this).closest('.vis-filter').remove();
            });

            visfilter.attr('data-column', column);
            visfilter.find('.name').text($(this).data('name'));
            visfilter.find('.filter').text(' (' + filters[column] + ')');
            visfilter.appendTo(visfilters);
        }

    });

// todo добавляется в скрытыйе боди всех .filter-condition
function AddingCondition() {
    var condition = $(this).data('filter-name');
    var filters = $(".filter-condition");
    var filter = $('.filter-condition-hidden .input-group[data-filter-name="' + condition + '"]');

    var cloneFilter = filter.clone();
    cloneFilter.find(".delete").on('click', function () {
        $(this).closest('.input-group').remove();
    });
    cloneFilter.appendTo(filters);
}


$('.reset').click(function () {
    $('#vis-filters .vis-filter').remove();
});

