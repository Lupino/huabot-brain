LESS_SOURCE =
JSX_SOURCE = app/main.js
APP = public/static/js/main.js

all: $(APP)

$(APP): $(JSX_SOURCE)
	cat $(JSX_SOURCE) > comibed.js
	browserify -t [ reactify --es6 ] comibed.js | uglifyjs -m -r '$$' > $@


clean:
	rm -f $(APP)
