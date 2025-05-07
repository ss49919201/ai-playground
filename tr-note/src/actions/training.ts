"use server";

import { Exercise, Set, TrainingRecord } from "../types";
import {
  addTrainingRecord as dbAddTrainingRecord,
  getTrainingRecord as dbGetTrainingRecord,
  getAllTrainingRecords as dbGetAllTrainingRecords,
  deleteTrainingRecord as dbDeleteTrainingRecord,
  addExerciseToRecord as dbAddExerciseToRecord,
  addSetToExercise as dbAddSetToExercise
} from "../db/training";

// トレーニング記録を追加する
export async function addTrainingRecord(
  record: Omit<TrainingRecord, "id" | "createdAt" | "updatedAt">
): Promise<TrainingRecord> {
  return dbAddTrainingRecord(record);
}

// トレーニング記録を取得する
export async function getTrainingRecord(
  id: string
): Promise<TrainingRecord | undefined> {
  return dbGetTrainingRecord(id);
}

// 全てのトレーニング記録を取得する
export async function getAllTrainingRecords(): Promise<TrainingRecord[]> {
  return dbGetAllTrainingRecords();
}

// トレーニング記録を削除する
export async function deleteTrainingRecord(id: string): Promise<void> {
  return dbDeleteTrainingRecord(id);
}

// トレーニング記録に種目を追加する
export async function addExerciseToRecord(
  recordId: string,
  exercise: Omit<Exercise, "id">
): Promise<void> {
  return dbAddExerciseToRecord(recordId, exercise);
}

// 種目にセットを追加する
export async function addSetToExercise(
  recordId: string,
  exerciseId: string,
  set: Omit<Set, "id">
): Promise<void> {
  return dbAddSetToExercise(recordId, exerciseId, set);
}
