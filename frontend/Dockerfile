# --- Этап сборки (Build Stage) ---
FROM node:20-alpine as build-stage

WORKDIR /app

# Копируем package.json и устанавливаем зависимости
COPY package*.json ./
RUN npm install

# Копируем исходники и собираем проект (Vite + Tailwind)
COPY . .
RUN npm run build

# --- Этап запуска (Production Stage) ---
FROM nginx:stable-alpine as production-stage

# Копируем собранные файлы (папка dist) в папку Nginx
COPY --from=build-stage /app/dist /usr/share/nginx/html

# Открываем 80 порт внутри контейнера
EXPOSE 80

# Запускаем Nginx
CMD ["nginx", "-g", "daemon off;"]