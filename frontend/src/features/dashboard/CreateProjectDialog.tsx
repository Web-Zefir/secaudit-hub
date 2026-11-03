import { useState, type FormEvent } from "react";
import { Button } from "../../components/ui/Button";
import { Modal } from "../../components/ui/Dialog";
import { Field, SelectInput, TextArea, TextInput } from "../../components/ui/Field";
import styles from "./Dashboard.module.scss";

interface CreateProjectDialogProps {
  open: boolean;
  pending: boolean;
  error: string;
  onOpenChange: (open: boolean) => void;
  onSubmit: (input: { name: string; description: string; repoUrl: string; language: string }) => void;
}

export function CreateProjectDialog({ open, pending, error, onOpenChange, onSubmit }: CreateProjectDialogProps) {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [repoUrl, setRepoUrl] = useState("");
  const [language, setLanguage] = useState("Go");

  function submit(event: FormEvent) {
    event.preventDefault();
    onSubmit({ name, description, repoUrl, language });
  }

  return (
    <Modal open={open} onOpenChange={onOpenChange} title="Новый проект">
      <form className={styles.form} onSubmit={submit}>
        <Field label="Название">
          <TextInput value={name} onChange={(event) => setName(event.target.value)} required minLength={2} />
        </Field>
        <Field label="Описание">
          <TextArea value={description} onChange={(event) => setDescription(event.target.value)} />
        </Field>
        <Field label="Репозиторий">
          <TextInput value={repoUrl} onChange={(event) => setRepoUrl(event.target.value)} placeholder="https://github.com/org/repo" />
        </Field>
        <Field label="Язык">
          <SelectInput value={language} onChange={(event) => setLanguage(event.target.value)}>
            <option>Go</option>
            <option>TypeScript</option>
            <option>Python</option>
            <option>Mixed</option>
          </SelectInput>
        </Field>
        {error ? <div>{error}</div> : null}
        <Button variant="primary" type="submit" disabled={pending}>
          Создать проект
        </Button>
      </form>
    </Modal>
  );
}
