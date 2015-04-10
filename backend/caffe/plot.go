package caffe

import (
    "log"
    "io/ioutil"
)

var PLOT_ROOT = resoursesPath + "/plot"

func Plot(suffix string) (data []byte, err error) {

    if err = run(PLOT_ROOT + "/parse_log.sh", LOG_DIR + "/caffe.INFO"); err != nil {
        log.Printf("Error: %s\n", err)
        return
    }
    if err = run("gnuplot", PLOT_ROOT + "/plot_log.gnuplot." + suffix); err != nil {
        log.Printf("Error: %s\n", err)
        return
    }
    data, err = ioutil.ReadFile("/tmp/" + suffix + ".png")
    return
}
