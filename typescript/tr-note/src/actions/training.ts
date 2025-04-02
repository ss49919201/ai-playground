"use server";

import { Exercise, Set, TrainingRecord } from "../types";

// メモリ上でデータを保持するための変数
let trainingRecords: TrainingRecord[] = [];

// ユニークIDを生成する関数
const generateId = (): string => {
  return Math.random().toString(36).substring(2, 9);
};

// 現在の日時を ISO 文字列で取得する関数
const getCurrentISOString = (): string => {
  return new Date().toISOString();
};

// トレーニング記録を追加する
export async function addTrainingRecord(
  record: Omit<TrainingRecord, "id" | "createdAt" | "updatedAt">
): Promise<TrainingRecord> {
  const timestamp = getCurrentISOString();
  const newRecord: TrainingRecord = {
    ...record,
    id: generateId(),
    createdAt: timestamp,
    updatedAt: timestamp,
  };

  trainingRecords = [...trainingRecords, newRecord];
  console.log("トレーニング記録が追加されました:", newRecord);
  return newRecord;
}

// トレーニング記録を取得する
export async function getTrainingRecord(
  id: string
): Promise<TrainingRecord | undefined> {
  const record = trainingRecords.find((record) => record.id === id);
  console.log(`ID: ${id} のトレーニング記録:`, record);
  return record;
}

// 全てのトレーニング記録を取得する
export async function getAllTrainingRecords(): Promise<TrainingRecord[]> {
  return trainingRecords;
}

// トレーニング記録を削除する
export async function deleteTrainingRecord(id: string): Promise<void> {
  trainingRecords = trainingRecords.filter((record) => record.id !== id);
  console.log(`ID: ${id} のトレーニング記録が削除されました`);
}

// トレーニング記録に種目を追加する
export async function addExerciseToRecord(
  recordId: string,
  exercise: Omit<Exercise, "id">
): Promise<void> {
  trainingRecords = trainingRecords.map((record) => {
    if (record.id === recordId) {
      const newExercise: Exercise = {
        ...exercise,
        id: generateId(),
      };
      const updatedRecord = {
        ...record,
        exercises: [...record.exercises, newExercise],
        updatedAt: getCurrentISOString(),
      };
      console.log(
        `ID: ${recordId} のトレーニング記録に種目が追加されました:`,
        newExercise
      );
      return updatedRecord;
    }
    return record;
  });
}

// 種目にセットを追加する
export async function addSetToExercise(
  recordId: string,
  exerciseId: string,
  set: Omit<Set, "id">
): Promise<void> {
  trainingRecords = trainingRecords.map((record) => {
    if (record.id === recordId) {
      const updatedExercises = record.exercises.map((exercise) => {
        if (exercise.id === exerciseId) {
          const newSet: Set = {
            ...set,
            id: generateId(),
          };
          console.log(
            `ID: ${exerciseId} の種目にセットが追加されました:`,
            newSet
          );
          return {
            ...exercise,
            sets: [...exercise.sets, newSet],
          };
        }
        return exercise;
      });

      return {
        ...record,
        exercises: updatedExercises,
        updatedAt: getCurrentISOString(),
      };
    }
    return record;
  });
}
