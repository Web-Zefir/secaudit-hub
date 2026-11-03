import { useState } from "react";
import { Button } from "../../components/ui/Button";
import { Modal } from "../../components/ui/Dialog";
import { Field, TextArea, TextInput } from "../../components/ui/Field";
import { TabList, TabPanel, TabTrigger, TabsRoot } from "../../components/ui/Tabs";
import type { CodeFile, TargetSource } from "../../types/models";
import styles from "./Project.module.scss";

interface LaunchScanDialogProps {
  open: boolean;
  pending: boolean;
  error: string;
  onOpenChange: (open: boolean) => void;
  onStart: (input: { targetSource: TargetSource; repoUrl?: string; files?: CodeFile[] }) => void;
}

export function LaunchScanDialog({ open, pending, error, onOpenChange, onStart }: LaunchScanDialogProps) {
  const [name, setName] = useState("snippet.go");
  const [content, setContent] = useState("");
  const [files, setFiles] = useState<CodeFile[]>([]);
  const [repoUrl, setRepoUrl] = useState("");

  function addFile() {
    const trimmed = name.trim();
    if (!trimmed || !content.trim()) return;
    setFiles((current) => [...current.filter((file) => file.name !== trimmed), { name: trimmed, content }]);
    setContent("");
  }

  return (
    <Modal open={open} onOpenChange={onOpenChange} title="Запустить SAST-аудит">
      <TabsRoot defaultValue="demo">
        <TabList>
          <TabTrigger value="demo">Демо-стенд уязвимого кода</TabTrigger>
          <TabTrigger value="code">Вставить свой код / файлы</TabTrigger>
          <TabTrigger value="git">Git-репозиторий по URL</TabTrigger>
        </TabList>
        <TabPanel value="demo">
          <div className={styles.stack}>
            <p className={styles.note}>
              Проверяет встроенный набор на Go, TypeScript, Python, SQL и .env. Все 12 правил должны дать находки.
            </p>
            <Button variant="primary" disabled={pending} onClick={() => onStart({ targetSource: "DEMO_PRESET" })}>
              Запустить демо-скан
            </Button>
          </div>
        </TabPanel>
        <TabPanel value="code">
          <div className={styles.stack}>
            <Field label="Имя файла">
              <TextInput value={name} onChange={(event) => setName(event.target.value)} />
            </Field>
            <Field label="Содержимое">
              <TextArea value={content} onChange={(event) => setContent(event.target.value)} spellCheck={false} />
            </Field>
            <Button onClick={addFile}>Добавить файл</Button>
            <div className={styles.fileList}>
              {files.map((file) => (
                <div key={file.name} className={styles.file}>
                  <span>{file.name}</span>
                  <button type="button" onClick={() => setFiles((current) => current.filter((item) => item.name !== file.name))}>
                    Убрать
                  </button>
                </div>
              ))}
            </div>
            <Button variant="primary" disabled={pending || files.length === 0} onClick={() => onStart({ targetSource: "CUSTOM_CODE", files })}>
              Сканировать файлы
            </Button>
          </div>
        </TabPanel>
        <TabPanel value="git">
          <div className={styles.stack}>
            <Field label="Публичный HTTPS URL">
              <TextInput value={repoUrl} onChange={(event) => setRepoUrl(event.target.value)} placeholder="https://github.com/org/repo.git" />
            </Field>
            <p className={styles.note}>Клонирование shallow (`git clone --depth 1`) во временную директорию. Локальные и приватные адреса отклоняются.</p>
            <Button variant="primary" disabled={pending || repoUrl.trim() === ""} onClick={() => onStart({ targetSource: "GIT_REPO", repoUrl })}>
              Клонировать и проверить
            </Button>
          </div>
        </TabPanel>
      </TabsRoot>
      {error ? <p className={styles.error}>{error}</p> : null}
    </Modal>
  );
}
