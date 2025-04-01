"use client";

import { createContext, useContext, useState, ReactNode } from "react";
import { TrainingRecord, Exercise, Set } from "../types";

// ユニークIDを生成する関数
const generateId = (): string => {
  return Math.random().toString(36).substring(2, 9);
};

// 現在の日時を ISO 文字列で取得する関数
const getCurrentISOString = (): string => {
  return new Date().toISOString();
};

interface TrainingContextType {
  trainingRecords: TrainingRecord[];
  addTrainingRecord: (
    record: Omit<TrainingRecord, "id" | "createdAt" | "updatedAt">
  ) => TrainingRecord;
  getTrainingRecord: (id: string) => TrainingRecord | undefined;
  deleteTrainingRecord: (id: string) => void;
  addExerciseToRecord: (
    recordId: string,
    exercise: Omit<Exercise, "id">
  ) => void;
  addSetToExercise: (
    recordId: string,
    exerciseId: string,
    set: Omit<Set, "id">
  ) => void;
}

const TrainingContext = createContext<TrainingContextType | undefined>(
  undefined
);

export const TrainingProvider = ({ children }: { children: ReactNode }) => {
  const [trainingRecords, setTrainingRecords] = useState<TrainingRecord[]>([]);

  // トレーニング記録を追加する
  const addTrainingRecord = (
    record: Omit<TrainingRecord, "id" | "createdAt" | "updatedAt">
  ): TrainingRecord => {
    const timestamp = getCurrentISOString();
    const newRecord: TrainingRecord = {
      ...record,
      id: generateId(),
      createdAt: timestamp,
      updatedAt: timestamp,
    };

    setTrainingRecords((prev) => [...prev, newRecord]);
    console.log("トレーニング記録が追加されました:", newRecord);
    return newRecord;
  };

  // トレーニング記録を取得する
  const getTrainingRecord = (id: string): TrainingRecord | undefined => {
    const record = trainingRecords.find((record) => record.id === id);
    console.log(`ID: ${id} のトレーニング記録:`, record);
    return record;
  };

  // トレーニング記録を削除する
  const deleteTrainingRecord = (id: string): void => {
    setTrainingRecords((prev) => prev.filter((record) => record.id !== id));
    console.log(`ID: ${id} のトレーニング記録が削除されました`);
  };

  // トレーニング記録に種目を追加する
  const addExerciseToRecord = (
    recordId: string,
    exercise: Omit<Exercise, "id">
  ): void => {
    setTrainingRecords((prev) =>
      prev.map((record) => {
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
      })
    );
  };

  // 種目にセットを追加する
  const addSetToExercise = (
    recordId: string,
    exerciseId: string,
    set: Omit<Set, "id">
  ): void => {
    setTrainingRecords((prev) =>
      prev.map((record) => {
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
      })
    );
  };

  const value = {
    trainingRecords,
    addTrainingRecord,
    getTrainingRecord,
    deleteTrainingRecord,
    addExerciseToRecord,
    addSetToExercise,
  };

  return (
    <TrainingContext.Provider value={value}>
      {children}
    </TrainingContext.Provider>
  );
};

export const useTraining = (): TrainingContextType => {
  const context = useContext(TrainingContext);
  if (context === undefined) {
    throw new Error("useTraining must be used within a TrainingProvider");
  }
  return context;
};
