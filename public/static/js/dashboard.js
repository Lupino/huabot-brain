"use strict"
var app = {};

app.$sidebar = $(".sidebar");
app.sidebarActive = ".sidebar ul li.active";
app.$container = $("#container");
app.currentPage = "#container .pagination .active";
app.$modal = $(".modal");
app.$modalBody = $(".modal .modal-body");
app.$modalBtns = $(".modal .modal-footer button");
app.target = "all";
app.modalReady = false;
app.datasets = [];

var hash = window.location.hash;
var m = /test|train|all/.exec(hash);
if (m) {
  app.target = m[0];
}

var hash_max, hash_limit;
var m = /max=(\d+)/.exec(hash);
if (m) {
  hash_max = m[1];
}

var m = /limit=(\d+)/.exec(hash);
if (m) {
  hash_limit = m[1];
}

var m = /test|train|all/.exec(hash);
if (m) {
  app.target = m[0];
}

app.$sidebar.click(function(evt) {
  var $elem = $(evt.target);
  if ($elem.prop("tagName") !== "A") {
    return;
  }
  if ($elem.parent().hasClass("active")) {
    return;
  }
  updateMenu($elem);
});

function updateMenu($elem) {
  $(app.sidebarActive).removeClass("active");
  $elem.parent().addClass("active");
  var target = $elem.attr('data-target');
  app.target = target;
  app.datasets = [];
  window.scroll(0, 0);
  loadData();
}

app.$container.click(function(evt) {
  var $elem = $(evt.target);
  if ($elem.prop("tagName") === "IMG" || $elem.attr('class') === 'tag') {
    $elem = $elem.parent();
    $elem = $elem.parent();
    var file_id = $elem.attr("data-id");
    $.get("/api/datasets/" + file_id, function(file) {
      var html = template(file);
      app.$modalBody.html(html);
      app.$modal.attr("data-id", file_id);
      app.$modal.modal();
    });
    return;
  }

  if ($elem.prop("tagName") !== "button") {
    $elem = $elem.parent();
    var max = $elem.attr('max');
    loadData(max);
    return;
  }

  if ($elem.prop("tagName") !== "A") {
    return;
  }
  if ($elem.parent().hasClass("active")) {
    return;
  }
  $(app.currentPage).removeClass("active");
  $elem.parent().addClass("active");
  var start = $elem.attr("data-start");
  loadData(app.target, start);
});


function loadData(max, limit, callback) {
  max = max || hash_max || "";// || Number.MAX_VALUE;
  limit = limit || hash_limit || 100;
  window.location.hash = "#" + app.target + '/max=' + max + '/limit=' + limit;
  $.get("/api/datasets/?max=" + max + "&limit=" + limit + "&data_type=" + app.target, function(data) {
    data.limit = limit;
    app.datasets = data.datasets = app.datasets.concat(data.datasets);
    data.datasets = app.datasets;
    var html = template(data);
    app.$container.html(html);
    $("#waterfall").waterfall({
      selector: ".img",
      // isResizable: true
    });
    if (callback) {
      callback();
    }
  });
}


app.$modalBtns.click(function(evt) {
  var $elem = $(evt.target);
  var file_id = app.$modal.attr("data-id");
  $.post("/api/datasets/" + file_id + "/", function(file) {
    app.$modal.modal("hide");
  });
});

updateMenu($('a[data-target=' + app.target + ']'));
