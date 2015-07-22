/**
 * Created by Vincent on 08/06/15.
 */

function Http() {
    this.Get = function (url) {
        return $.ajax({
            method: 'GET',
            url: url,
            dataType: 'json'
        });
    };

    this.Post = function (url, data) {
        return $.ajax({
            method: 'POST',
            url: url,
            data: JSON.stringify(data),
            dataType: 'json',
            contentType: 'application/json;charset=UTF-8'
        });
    };

    this.Delete = function (url, data) {
        return $.ajax({
            method: 'DELETE',
            url: url,
            data: JSON.stringify(data),
            dataType: 'json',
            contentType: 'application/json;charset=UTF-8'
        });
    };

    this.Put = function (url, data) {
        return $.ajax({
            method: 'PUT',
            url: url,
            data: JSON.stringify(data),
            dataType: 'json',
            contentType: 'application/json;charset=UTF-8'
        });
    };
}