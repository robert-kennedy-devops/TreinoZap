export interface Trainer {
  id: string;
  name: string;
  email: string;
  phone: string;
  role: string;
  status: string;
  created_at: string;
}

export interface Client {
  id: string;
  trainer_id: string;
  name: string;
  phone: string;
  status: string;
  goal?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export interface Exercise {
  id: string;
  trainer_id: string;
  name: string;
  muscle_group?: string;
  equipment?: string;
  video_url?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export interface WorkoutExercise {
  id: string;
  section_id: string;
  exercise_id?: string;
  exercise_name: string;
  sets?: string;
  reps?: string;
  rest_seconds?: number;
  load_note?: string;
  technique_note?: string;
  video_url?: string;
  order_index: number;
}

export interface WorkoutSection {
  id: string;
  workout_id: string;
  name: string;
  description?: string;
  order_index: number;
  exercises?: WorkoutExercise[];
}

export interface Workout {
  id: string;
  trainer_id: string;
  client_id: string;
  name: string;
  status: string;
  starts_at?: string;
  ends_at?: string;
  sections?: WorkoutSection[];
  created_at: string;
  updated_at: string;
}

export interface Message {
  id: string;
  trainer_id?: string;
  client_id?: string;
  direction: 'inbound' | 'outbound';
  phone: string;
  message: string;
  command?: string;
  status?: string;
  created_at: string;
}

export interface WhatsAppStatus {
  connected: boolean;
  phone?: string;
  jid?: string;
  last_connected?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface ApiError {
  error: {
    message: string;
    code: string;
  };
}
