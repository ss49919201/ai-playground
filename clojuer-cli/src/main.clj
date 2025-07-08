(ns main
  (:require [clojure.test :refer :all]
            [clojure.string]
            [ring.adapter.jetty :refer [run-jetty]]
            [compojure.core :refer [defroutes GET POST]]
            [compojure.route :as route]
            [ring.middleware.json :refer [wrap-json-response]]
            [ring.util.response :refer [response status]]
            [ring.middleware.json :refer [wrap-json-body]]))

(def users-db (atom {}))

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

(defn quicksort [coll]
  ;; 空の配列または要素が1つの場合はそのまま返す
  (if (<= (count coll) 1)
    coll
    ;; ピボットを最初の要素とする
    (let [pivot (first coll)
          ;; 残りの要素を取得
          rest-coll (rest coll)
          ;; ピボットより小さい要素を抽出
          smaller (filter #(< % pivot) rest-coll)
          ;; ピボットより大きい要素を抽出
          larger (filter #(>= % pivot) rest-coll)]
      ;; 小さい要素をソート + ピボット + 大きい要素をソート
      (concat (quicksort smaller) [pivot] (quicksort larger)))))

(defn quicksort-iterative [coll]
  ;; 簡単な実装: 実際には再帰を使わずに実装するのは複雑
  ;; ここでは参考として単純な実装を示す
  (if (<= (count coll) 1)
    coll
    ;; ループで実装（結果を段階的に構築）
    (loop [to-process [coll]
           completed []]
      ;; 処理待ちがなくなったら結果を返す
      (if (empty? to-process)
        (flatten completed)
        ;; 次の配列を処理
        (let [current (first to-process)
              remaining (rest to-process)]
          ;; 1要素以下なら完了リストに追加
          (if (<= (count current) 1)
            (recur remaining (conj completed current))
            ;; 分割して処理待ちに追加
            (let [pivot (first current)
                  rest-coll (rest current)
                  smaller (filter #(< % pivot) rest-coll)
                  larger (filter #(>= % pivot) rest-coll)]
              ;; 小さい要素、ピボット、大きい要素の順で処理
              (recur (concat [smaller [pivot] larger] remaining) completed))))))))

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

(deftest quicksort-test
  (testing "Quicksort function"
    (is (= [] (quicksort [])))
    (is (= [1] (quicksort [1])))
    (is (= [1 2 3] (quicksort [3 1 2])))
    (is (= [1 2 3 5 8 9] (quicksort [5 2 8 1 9 3]))))
  (testing "Quicksort iterative function"
    (is (= [] (quicksort-iterative [])))
    (is (= [1] (quicksort-iterative [1])))
    (is (= [1 2 3] (quicksort-iterative [3 1 2])))
    (is (= [1 2 3 5 8 9] (quicksort-iterative [5 2 8 1 9 3])))))

(defn health-handler [request]
  (response {:msg "ok"}))

(defn get-user-handler [request]
  (let [id (get-in request [:route-params :id])
        user (get @users-db id)]
    (if user
      (response user)
      (-> (response {:error "User not found"})
          (status 404)))))

(defn create-user-handler [request]
  (let [user-data (:body request)
        id (str (java.util.UUID/randomUUID))]
    (if user-data
      (do
        (swap! users-db assoc id user-data)
        (-> (response (assoc user-data :id id))
            (status 201)))
      (-> (response {:error "Invalid user data"})
          (status 400)))))

(deftest health-handler-test
  (testing "Health endpoint handler"
    (let [response (health-handler {})]
      (is (= {:msg "ok"} (:body response)))
      (is (= 200 (:status response))))))

(deftest user-handlers-test
  (testing "Get user handler"
    (reset! users-db {"test-id" {:name "Test User" :email "test@example.com"}})
    (let [response (get-user-handler {:route-params {:id "test-id"}})]
      (is (= {:name "Test User" :email "test@example.com"} (:body response)))
      (is (= 200 (:status response))))
    (let [response (get-user-handler {:route-params {:id "nonexistent"}})]
      (is (= {:error "User not found"} (:body response)))
      (is (= 404 (:status response)))))
  (testing "Create user handler"
    (reset! users-db {})
    (let [user-data {:name "New User" :email "new@example.com"}
          response (create-user-handler {:body user-data})]
      (is (= "New User" (get-in response [:body :name])))
      (is (= "new@example.com" (get-in response [:body :email])))
      (is (string? (get-in response [:body :id])))
      (is (= 201 (:status response)))
      (is (= 1 (count @users-db))))
    (let [response (create-user-handler {:body nil})]
      (is (= {:error "Invalid user data"} (:body response)))
      (is (= 400 (:status response))))))

(defroutes app-routes
  (GET "/health" [] health-handler)
  (GET "/users/:id" [id] get-user-handler)
  (POST "/users/" [] create-user-handler)
  (route/not-found "Not Found"))

(def app
  (-> app-routes
      (wrap-json-body {:keywords? true})
      wrap-json-response))

(defn -main [& args]
  (println "Starting server on port 8080...")
  (run-jetty app {:port 8080 :join? false}))
