import { Link } from "@tanstack/react-router";
import styles from "./NotFoundPage.module.scss";

export function NotFoundPage() {
  return (
    <div className={styles.page}>
      <h1>Страница не найдена</h1>
      <Link to="/">Вернуться к обзору</Link>
    </div>
  );
}
