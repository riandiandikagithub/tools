import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function parseHumanSize(str: string): number {
  const units: Record<string, number> = {
    K: 1024,
    M: 1024 ** 2,
    G: 1024 ** 3,
  };

  const match = str.match(/([\d.]+)([KMG]?)/);
  if (!match) return 0;

  const value = parseFloat(match[1]);
  const unit = match[2];

  return value * (units[unit] || 1);
}

export function calculateReplicationLag(
  masterReplOffset?: number,
  replicaOffset?: number,
): number {
  if (
    masterReplOffset === undefined ||
    replicaOffset === undefined ||
    isNaN(masterReplOffset) ||
    isNaN(replicaOffset)
  ) {
    return 0;
  }

  const lag = masterReplOffset - replicaOffset;
  return lag < 0 ? 0 : lag;
}
