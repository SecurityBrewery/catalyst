export interface Problem {
  title: string;
  detail: string;
}

enum AlertType {
  success = "success",
  info = "info",
  warning = "warning",
  error = "error",
}

export interface Alert {
  name: string;
  detail: string;
  type: AlertType;
}
