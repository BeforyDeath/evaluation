/// <reference path="../typings/jquery/jquery.d.ts" />

class Records {
    public data = {};

    load = (url:string) => {
        let self = this;
        $.ajax({
            url: url,
            dataType: 'json',
            async: false,
            success: function (data) {
                self.data = data;
            }
        });
    };

    getSecret = () => {
        let data = false;
        $.ajax({
            url: "api/secret",
            dataType: 'json',
            async: false,
            success: function (_data) {
                data = _data
            },
        });
        return data
    };
    
}