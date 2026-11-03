import { forwardRef, type ButtonHTMLAttributes } from "react";
import { cx } from "../../lib/format";
import styles from "./Button.module.scss";

type Variant = "primary" | "ghost" | "danger";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant;
  wide?: boolean;
}

const variantClass: Record<Variant, string> = {
  primary: styles.primary,
  ghost: styles.ghost,
  danger: styles.danger,
};

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(function Button(
  { variant = "ghost", wide = false, className, type = "button", ...props },
  ref,
) {
  return (
    <button
      ref={ref}
      type={type}
      className={cx(styles.button, variantClass[variant], wide && styles.wide, className)}
      {...props}
    />
  );
});
