package main

import (
    "flag"
)

var (
    API = "/api"
    GEARMAND = flag.String("gearmand", "127.0.0.1:4730", "The Gearmand server.")
)
