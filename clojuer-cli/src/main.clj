(ns main
  (:require [clojure.test :refer :all]
            [clojure.string]))

(defn fizzbuzz [n]
  (cond
    (= (mod n 15) 0) "FizzBuzz"
    (= (mod n 3) 0) "Fizz"
    (= (mod n 5) 0) "Buzz"
    :else (str n)))

(defn run-fizzbuzz [max]
  (doseq [i (range 1 (inc max))]
    (println (fizzbuzz i))))

(defn prime? [n]
  (cond
    (< n 2) false
    (= n 2) true
    (even? n) false
    :else (let [sqrt-n (int (Math/sqrt n))]
            (not-any? #(zero? (mod n %)) (range 3 (inc sqrt-n) 2)))))

(defn filter-primes [nums]
  (filter prime? nums))

(defn lowercase? [s]
  (= s (clojure.string/lower-case s)))

(defn not-all-lowercase? [strings]
  (not-every? lowercase? strings))

(deftest fizzbuzz-test
  (testing "FizzBuzz function"
    (is (= "1" (fizzbuzz 1)))
    (is (= "2" (fizzbuzz 2)))
    (is (= "Fizz" (fizzbuzz 3)))
    (is (= "4" (fizzbuzz 4)))
    (is (= "Buzz" (fizzbuzz 5)))
    (is (= "Fizz" (fizzbuzz 6)))
    (is (= "Fizz" (fizzbuzz 9)))
    (is (= "Buzz" (fizzbuzz 10)))
    (is (= "FizzBuzz" (fizzbuzz 15)))
    (is (= "FizzBuzz" (fizzbuzz 30)))))

(deftest prime-filter-test
  (testing "Prime number filter function"
    (is (= [2 3] (filter-primes [1 2 3 4])))
    (is (= [] (filter-primes [1 4 6 8])))
    (is (= [2 3 5 7] (filter-primes [1 2 3 4 5 6 7 8 9 10])))))

(deftest not-all-lowercase-test
  (testing "Not all lowercase function"
    (is (= true (not-all-lowercase? ["hello" "World"])))
    (is (= false (not-all-lowercase? ["hello" "world"])))
    (is (= true (not-all-lowercase? ["Hello" "WORLD"])))
    (is (= false (not-all-lowercase? [])))))

(defn -main [& args]
  (run-fizzbuzz 100))
