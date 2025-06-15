# Estágio de build
FROM node:20-alpine AS build-stage

WORKDIR /app

# Copia o package.json e o lock file
COPY package*.json ./

# Instala as dependências
RUN npm install

# Copia o restante dos arquivos do frontend
COPY . .

# Compila a aplicação para produção
RUN npm run build

# Estágio de produção
FROM nginx:stable-alpine AS production-stage

# Copia os arquivos compilados do estágio de build para o diretório padrão do Nginx
COPY --from=build-stage /app/dist /usr/share/nginx/html

# Remove a configuração padrão do Nginx
RUN rm /etc/nginx/conf.d/default.conf

# Copia o nosso arquivo de configuração personalizado
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expõe a porta 80 para acesso ao frontend
EXPOSE 80

# Comando para iniciar o Nginx
CMD ["nginx", "-g", "daemon off;"]