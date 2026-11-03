import * as Tabs from "@radix-ui/react-tabs";
import type { ReactNode } from "react";
import styles from "./Tabs.module.scss";

export function TabList({ children }: { children: ReactNode }) {
  return <Tabs.List className={styles.list}>{children}</Tabs.List>;
}

export function TabTrigger({ value, children }: { value: string; children: string }) {
  return (
    <Tabs.Trigger className={styles.trigger} value={value}>
      {children}
    </Tabs.Trigger>
  );
}

export function TabPanel({ value, children }: { value: string; children: ReactNode }) {
  return (
    <Tabs.Content className={styles.panel} value={value}>
      {children}
    </Tabs.Content>
  );
}

export const TabsRoot = Tabs.Root;
