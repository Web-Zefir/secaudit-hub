import { Link, useNavigate } from "@tanstack/react-router";
import { useEffect, useState, type FormEvent } from "react";
import { errorMessage } from "../../api/client";
import { useAuth } from "../../context/AuthContext";
import { NeutralBadge, RoadmapBadge } from "../../components/ui/Badge";
import { Button } from "../../components/ui/Button";
import { Field, TextInput } from "../../components/ui/Field";
import styles from "./AuthScreen.module.scss";

interface AuthScreenProps {
  mode: "login" | "register";
}

export function AuthScreen({ mode }: AuthScreenProps) {
  const { user, loading, login, register, loginDemo } = useAuth();
  const navigate = useNavigate();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [pending, setPending] = useState(false);

  useEffect(() => {
    if (!loading && user) {
      void navigate({ to: "/" });
    }
  }, [loading, user, navigate]);

  async function submit(event: FormEvent) {
    event.preventDefault();
    setPending(true);
    setError("");
    try {
      if (mode === "login") {
        await login(email, password);
      } else {
        await register(name, email, password);
      }
      void navigate({ to: "/" });
    } catch (err) {
      setError(errorMessage(err));
    } finally {
      setPending(false);
    }
  }

  async function demo() {
    setPending(true);
    setError("");
    try {
      await loginDemo();
      void navigate({ to: "/" });
    } catch (err) {
      setError(errorMessage(err));
    } finally {
      setPending(false);
    }
  }

  return (
    <div className={styles.page}>
      <aside className={styles.aside}>
        <div className={styles.kicker}>SECAUDIT HUB</div>
        <h1 className={styles.title}>Непрерывный контроль AppSec</h1>
        <p className={styles.lead}>Статический анализ исходников, отчёт по CWE и OWASP, разбор патча по каждой находке.</p>
        <div className={styles.modules}>
          <div className={styles.row}>
            <span>SAST</span>
            <NeutralBadge>Available</NeutralBadge>
          </div>
          <div className={styles.row}>
            <span>DAST</span>
            <RoadmapBadge />
          </div>
          <div className={styles.row}>
            <span>IaC</span>
            <RoadmapBadge />
          </div>
          <div className={styles.row}>
            <span>DB Security</span>
            <RoadmapBadge />
          </div>
        </div>
      </aside>
      <section className={styles.card}>
        <h2 className={styles.heading}>{mode === "login" ? "Вход" : "Регистрация"}</h2>
        <form className={styles.form} onSubmit={(event) => void submit(event)}>
          {mode === "register" ? (
            <Field label="Имя">
              <TextInput value={name} onChange={(event) => setName(event.target.value)} required minLength={2} />
            </Field>
          ) : null}
          <Field label="Email">
            <TextInput type="email" value={email} onChange={(event) => setEmail(event.target.value)} required autoComplete="username" />
          </Field>
          <Field label="Пароль">
            <TextInput
              type="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              required
              minLength={8}
              autoComplete={mode === "login" ? "current-password" : "new-password"}
            />
          </Field>
          {error ? <div className={styles.error}>{error}</div> : null}
          <Button variant="primary" type="submit" wide disabled={pending}>
            {mode === "login" ? "Войти" : "Создать аккаунт"}
          </Button>
        </form>
        <Button wide disabled={pending} onClick={() => void demo()}>
          Быстрый вход демо-аудитора
        </Button>
        <div className={styles.switch}>
          {mode === "login" ? (
            <>
              Нет аккаунта? <Link to="/register">Регистрация</Link>
            </>
          ) : (
            <>
              Уже есть аккаунт? <Link to="/login">Войти</Link>
            </>
          )}
        </div>
      </section>
    </div>
  );
}
