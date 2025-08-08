import styles from "./page.module.css";

const todo = [
  {
    id: 1,
    title: "プロジェクトの企画書を作成",
    body: "新しいアプリケーションの企画書を作成し、要件定義を明確にする",
    completed: false,
  },
  {
    id: 2,
    title: "デザインのレビュー",
    body: "UIデザインのレビューを行い、チームからフィードバックを収集する",
    completed: true,
  },
  {
    id: 3,
    title: "データベース設計",
    body: "アプリケーションで使用するデータベーススキーマを設計する",
    completed: false,
  },
  {
    id: 4,
    title: "テストケースの作成",
    body: "ユニットテストと統合テストのテストケースを作成する",
    completed: false,
  },
  {
    id: 5,
    title: "ドキュメント更新",
    body: "README.mdとAPI仕様書を最新の状態に更新する",
    completed: true,
  },
];

export default function Home() {
  return (
    <div className={styles.page}>
      <main className={styles.main}>
        <h1>TODO一覧</h1>
        <div className={styles.todoList}>
          {todo.map((item) => (
            <div key={item.id} className={styles.todoItem}>
              <div className={styles.todoHeader}>
                <h2 className={`${styles.todoTitle} ${item.completed ? styles.completed : ""}`}>
                  {item.title}
                </h2>
              </div>
              <p className={styles.todoBody}>{item.body}</p>
            </div>
          ))}
        </div>
      </main>
    </div>
  );
}
