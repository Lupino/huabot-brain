package caffe

import (
    "log"
    "io/ioutil"
    "github.com/Lupino/huabot-brain/config"
)

func Plot(suffix string) (data []byte, err error) {

    if err = run(config.PLOT_ROOT + "/parse_log.sh",
                 config.LOG_DIR + "/caffe.INFO"); err != nil {
        log.Printf("Error: %s\n", err)
        return
    }
    if err = run("gnuplot", config.PLOT_ROOT + "/plot_log.gnuplot." + suffix); err != nil {
        log.Printf("Error: %s\n", err)
        return
    }
    data, err = ioutil.ReadFile("/tmp/" + suffix + ".png")
    return
}
