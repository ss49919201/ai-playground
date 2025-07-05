(ns main
  (:gen-class)
  (:require [greeter]))

(defn -main [& args]
  (println (greeter/greeting)))
