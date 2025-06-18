<template>
  <Form @submit="handleSubmit" :validation-schema="schema" v-slot="{ errors }">
    <div class="form-group">
      <label for="name">Nome:</label>
      <Field name="name" type="text" id="name" v-model="formData.name" />
      <span class="error-message">{{ errors.name }}</span>
    </div>
    <div class="form-group">
      <label for="email">Email:</label>
      <Field name="email" type="email" id="email" v-model="formData.email" />
      <span class="error-message">{{ errors.email }}</span>
    </div>
    <div class="form-group" v-if="!isEditing">
      <label for="password">Senha:</label>
      <Field name="password" type="password" id="password" v-model="formData.password" />
      <span class="error-message">{{ errors.password }}</span>
    </div>
    <div class="form-actions">
      <button type="submit" class="btn-primary">Salvar</button>
      <button type="button" @click="close">Cancelar</button>
    </div>
  </Form>
  </template>

<script setup>
import { ref, watch, computed } from 'vue';
import { Form, Field } from 'vee-validate';
import * as yup from 'yup';

const props = defineProps({
  userToEdit: {
    type: Object,
    default: null,
  },
});

const emit = defineEmits(['submit', 'close']);


const isEditing = ref(false);
const formData = ref({});

const schema = yup.object({
  name: yup.string().required('O nome é obrigatório.'),
  email: yup.string().required('O email é obrigatório.').email('Formato de email inválido.'),
  password: yup.string().when([], {
    is: () => !isEditing.value,
    then: (schema) => schema.required('A senha é obrigatória.').min(8, 'A senha deve ter no mínimo 8 caracteres.'),
    otherwise: (schema) => schema.notRequired(),
  }),
});

const formTitle = computed(() => (props.userToEdit ? 'Editar Usuário' : 'Criar Novo Usuário'));

// Observa mudanças no prop para preencher o formulário para edição
watch(() => props.userToEdit, (newUser) => {
  if (newUser && newUser.id) { // Check for newUser.id to confirm it's an existing user
    formData.value = {
      id: newUser.id,
      name: newUser.name,
      email: newUser.email,
      password: '', // Password field is not for editing existing user's password here
    };
    isEditing.value = true;
  } else {
    // Novo usuário ou formulário limpo
    formData.value = {
      id: null,
      name: '',
      email: '',
      password: '', // Inicializa o campo de senha para novos usuários
    };
    isEditing.value = false;
  }
}, { immediate: true });

const handleSubmit = () => {
  emit('submit', { ...formData.value });
};
</script>

<style scoped>
.form-container {
  border: 1px solid #ccc;
  padding: 20px;
  margin-top: 20px;
  border-radius: 8px;
}
div {
  margin-bottom: 10px;
}
label {
  display: block;
  margin-bottom: 5px;
}
input {
  width: 100%;
  padding: 8px;
  box-sizing: border-box;
}
.actions {
  margin-top: 15px;
}
.error-message { color: red; font-size: 0.8em; }
</style>