function template(locals) {
var jade_debug = [{ lineno: 1, filename: "container.jade" }];
try {
var buf = [];
var jade_mixins = {};
var jade_interp;
;var locals_for_with = (locals || {});(function (datasets, undefined, dataset) {
jade_debug.unshift({ lineno: 0, filename: "container.jade" });
jade_debug.unshift({ lineno: 1, filename: "container.jade" });
if (datasets) {
{
jade_debug.unshift({ lineno: 2, filename: "container.jade" });
jade_debug.unshift({ lineno: 2, filename: "container.jade" });
buf.push("<div id=\"waterfall\">");
jade_debug.unshift({ lineno: undefined, filename: jade_debug[0].filename });
jade_debug.unshift({ lineno: 3, filename: "container.jade" });
// iterate datasets
;(function(){
  var $$obj = datasets;
  if ('number' == typeof $$obj.length) {

    for (var $index = 0, $$l = $$obj.length; $index < $$l; $index++) {
      var val = $$obj[$index];

jade_debug.unshift({ lineno: 3, filename: "container.jade" });
jade_debug.unshift({ lineno: 4, filename: "container.jade" });
buf.push("<div" + (jade.attr("data-id", val.dataset_id, true, false)) + " class=\"img\">");
jade_debug.unshift({ lineno: undefined, filename: jade_debug[0].filename });
jade_debug.unshift({ lineno: 5, filename: "container.jade" });
buf.push("<div class=\"dataset\">");
jade_debug.unshift({ lineno: undefined, filename: jade_debug[0].filename });
jade_debug.unshift({ lineno: 6, filename: "container.jade" });
buf.push("<img" + (jade.attr("src", "/upload/" + val.file.key, true, false)) + " width=\"192\"" + (jade.attr("height", 192/val.file.width*val.file.height, true, false)) + "/>");
jade_debug.shift();
jade_debug.unshift({ lineno: 7, filename: "container.jade" });
buf.push("<div class=\"tag\">");
jade_debug.unshift({ lineno: undefined, filename: jade_debug[0].filename });
jade_debug.unshift({ lineno: 7, filename: jade_debug[0].filename });
buf.push("" + (jade.escape((jade_interp = val.tag.name) == null ? '' : jade_interp)) + "");
jade_debug.shift();
jade_debug.shift();
buf.push("</div>");
jade_debug.shift();
jade_debug.shift();
buf.push("</div>");
jade_debug.shift();
jade_debug.shift();
buf.push("</div>");
jade_debug.shift();
jade_debug.shift();
    }

  } else {
    var $$l = 0;
    for (var $index in $$obj) {
      $$l++;      var val = $$obj[$index];

jade_debug.unshift({ lineno: 3, filename: "container.jade" });
jade_debug.unshift({ lineno: 4, filename: "container.jade" });
buf.push("<div" + (jade.attr("data-id", val.dataset_id, true, false)) + " class=\"img\">");
jade_debug.unshift({ lineno: undefined, filename: jade_debug[0].filename });
jade_debug.unshift({ lineno: 5, filename: "container.jade" });
buf.push("<div class=\"dataset\">");
jade_debug.unshift({ lineno: undefined, filename: jade_debug[0].filename });
jade_debug.unshift({ lineno: 6, filename: "container.jade" });
buf.push("<img" + (jade.attr("src", "/upload/" + val.file.key, true, false)) + " width=\"192\"" + (jade.attr("height", 192/val.file.width*val.file.height, true, false)) + "/>");
jade_debug.shift();
jade_debug.unshift({ lineno: 7, filename: "container.jade" });
buf.push("<div class=\"tag\">");
jade_debug.unshift({ lineno: undefined, filename: jade_debug[0].filename });
jade_debug.unshift({ lineno: 7, filename: jade_debug[0].filename });
buf.push("" + (jade.escape((jade_interp = val.tag.name) == null ? '' : jade_interp)) + "");
jade_debug.shift();
jade_debug.shift();
buf.push("</div>");
jade_debug.shift();
jade_debug.shift();
buf.push("</div>");
jade_debug.shift();
jade_debug.shift();
buf.push("</div>");
jade_debug.shift();
jade_debug.shift();
    }

  }
}).call(this);

jade_debug.shift();
jade_debug.shift();
buf.push("</div>");
jade_debug.shift();
jade_debug.unshift({ lineno: 8, filename: "container.jade" });
if (datasets.length > 0) {
{
jade_debug.unshift({ lineno: 9, filename: "container.jade" });
jade_debug.unshift({ lineno: 9, filename: "container.jade" });
buf.push("<div" + (jade.attr("max", datasets[datasets.length - 1].dataset_id, true, false)) + " class=\"load-more\">");
jade_debug.unshift({ lineno: undefined, filename: jade_debug[0].filename });
jade_debug.unshift({ lineno: 10, filename: "container.jade" });
buf.push("<button class=\"btn btn-lg btn-info\">");
jade_debug.unshift({ lineno: undefined, filename: jade_debug[0].filename });
jade_debug.unshift({ lineno: 10, filename: jade_debug[0].filename });
buf.push("加载更多...");
jade_debug.shift();
jade_debug.shift();
buf.push("</button>");
jade_debug.shift();
jade_debug.shift();
buf.push("</div>");
jade_debug.shift();
jade_debug.shift();
}
jade_debug.shift();
jade_debug.unshift({ lineno: 11, filename: "container.jade" });
}
jade_debug.shift();
jade_debug.shift();
}
jade_debug.shift();
jade_debug.unshift({ lineno: 12, filename: "container.jade" });
}
jade_debug.shift();
jade_debug.unshift({ lineno: 14, filename: "container.jade" });
if (dataset) {
{
jade_debug.unshift({ lineno: 15, filename: "container.jade" });
jade_debug.unshift({ lineno: 15, filename: "container.jade" });
buf.push("<div" + (jade.attr("data-id", dataset.dataset_id, true, false)) + " class=\"img\">");
jade_debug.unshift({ lineno: undefined, filename: jade_debug[0].filename });
jade_debug.unshift({ lineno: 16, filename: "container.jade" });
buf.push("<img" + (jade.attr("src", "/upload/" + dataset.file.key, true, false)) + " style=\"width: 100%;\"/>");
jade_debug.shift();
jade_debug.shift();
buf.push("</div>");
jade_debug.shift();
jade_debug.shift();
}
jade_debug.shift();
jade_debug.unshift({ lineno: 17, filename: "container.jade" });
}
jade_debug.shift();
jade_debug.shift();}("datasets" in locals_for_with?locals_for_with.datasets:typeof datasets!=="undefined"?datasets:undefined,"undefined" in locals_for_with?locals_for_with.undefined:typeof undefined!=="undefined"?undefined:undefined,"dataset" in locals_for_with?locals_for_with.dataset:typeof dataset!=="undefined"?dataset:undefined));;return buf.join("");
} catch (err) {
  jade.rethrow(err, jade_debug[0].filename, jade_debug[0].lineno, "- if (datasets) {\n  div#waterfall\n    each val in datasets\n      div.img(data-id=val.dataset_id)\n        div.dataset\n          img(src=\"/upload/\" + val.file.key,width=192,height=192/val.file.width*val.file.height)\n          div.tag #{val.tag.name}\n  - if (datasets.length > 0) {\n    div.load-more(max=datasets[datasets.length - 1].dataset_id)\n      button.btn.btn-lg.btn-info 加载更多...\n  - }\n- }\n\n- if (dataset) {\n  div.img(data-id=dataset.dataset_id)\n    img(src=\"/upload/\" + dataset.file.key,style=\"width: 100%;\")\n- }\n");
}
}