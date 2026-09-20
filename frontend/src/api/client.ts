import type {
  BackupProgress,
  Dataset,
  FileEntry,
  Node,
  UserInfo,
} from "./types";

import { GenerateRecoveryPhrase } from "../../wailsjs/go/main/App";

export async function unlock(recoveryPhrase: string): Promise<boolean> {
  void recoveryPhrase;

  throw new Error("Wails unlock method has not been connected yet.");
}

export async function getUserInfo(): Promise<UserInfo> {
  throw new Error("Wails getUserInfo method has not been connected yet.");
}

export async function listDatasets(): Promise<Dataset[]> {
  throw new Error("Wails listDatasets method has not been connected yet.");
}

export async function getDatasetFiles(datasetId: string): Promise<FileEntry[]> {
  void datasetId;

  throw new Error("Wails getDatasetFiles method has not been connected yet.");
}

export async function createDataset(label: string): Promise<Dataset> {
  void label;

  throw new Error("Wails createDataset method has not been connected yet.");
}

export async function startBackup(
  datasetId: string,
  paths: string[],
): Promise<void> {
  void datasetId;
  void paths;

  throw new Error("Wails startBackup method has not been connected yet.");
}

export async function restore(
  datasetId: string,
  paths: string[],
  destination: string,
): Promise<void> {
  void datasetId;
  void paths;
  void destination;

  throw new Error("Wails restore method has not been connected yet.");
}

export async function listNodes(): Promise<Node[]> {
  throw new Error("Wails listNodes method has not been connected yet.");
}

export async function addBootstrapNode(
  peerId: string,
  address: string,
): Promise<void> {
  void peerId;
  void address;

  throw new Error("Wails addBootstrapNode method has not been connected yet.");
}

export async function removeBootstrapNode(peerId: string): Promise<void> {
  void peerId;

  throw new Error(
    "Wails removeBootstrapNode method has not been connected yet.",
  );
}

export async function getBackupProgress(): Promise<BackupProgress> {
  throw new Error("Wails getBackupProgress method has not been connected yet.");
}

export async function generateRecoveryPhrase(): Promise<string> {
  return GenerateRecoveryPhrase();
}
