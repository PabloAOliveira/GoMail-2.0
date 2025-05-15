<template>
  <v-container>
    <v-form v-model="valid" ref="form" lazy="validation">
      <v-text-field
        label="Email"
        v-model="email.to"
        :rules="emailRules"
        required
      />
      <v-text-field
        label="Assunto"
        v-model="email.subject"
        required
        :rules="subjectRules"
      />

      <v-textarea
        label="Mensagem"
        v-model="email.body"
        required
        :rules="bodyRules"
      />

      <v-btn :disabled="!valid" @click="sendEmail" required color="primary">
        Enviar E-mail
      </v-btn>
    </v-form>

    <v-spacer class="mt-4"></v-spacer>

    <v-row>
      <template v-if="carregandoEmails">
        <v-col cols="12" md="4" v-for="n in 3" :key="n">
          <v-skeleton-loader
            class="mx-auto"
            max-width="400"
            type="image, heading, text"
          />
        </v-col>
      </template>

      <template v-else>
        <v-col
          v-for="message in emails.sort(
            (a, b) => new Date(b.created_at) - new Date(a.created_at)
          )"
          :key="message.id"
          cols="12"
          md="4"
        >
          <v-card class="mx-auto" hover>
            <v-card-item>
              <v-card-title>{{ message.subject }}</v-card-title>
              <v-card-subtitle>{{ message.to }}</v-card-subtitle>
            </v-card-item>

            <v-card-text>{{ message.body }}</v-card-text>

            <v-chip :color="statusColor(message.status)">
              Status: {{ statusLabel(message.status) }}
            </v-chip>
            <v-chip>Enviado em: {{ message.created_at }}</v-chip>
          </v-card>
        </v-col>
      </template>
    </v-row>
  </v-container>
</template>

<script setup>
import { onMounted, ref } from "vue";
import axios from "axios";

const form = ref(null);
const valid = ref(false);
const email = ref({
  to: "",
  subject: "",
  body: "",
});
const emails = ref([]);
const carregandoEmails = ref(false);

const emailRules = [
  (v) => !!v || "E-mail é obrigatório",
  (v) => /.+@.+\..+/.test(v) || "E-mail deve ser válido",
];

const subjectRules = [(v) => !!v || "Assunto é Obrigatório"];

const bodyRules = [(v) => !!v || "Mensagem é Obrigatória"];

onMounted(() => {
  fetchEmails();
});

const atualizaEmails = async () => {
  try {
  } catch (error) {}
};

const fetchEmails = async () => {
  try {
    const response = await axios.get("http://localhost:2000/emails");

    emails.value = response.data;
  } catch (error) {
    console.error(error);
    alert("Erro ao buscar E-mails");
  }
};

const statusColor = (status) => {
  const statusEnum = { sent: "green", enqueued: "orange", failed: "red" };

  return statusEnum[status];
};

const statusLabel = (status) => {
  switch (status) {
    case "sent":
      return "Enviado";
    case "enqueued":
      return "Pendente";
    case "failed":
      return "Falhou";
    default:
      return "Desconhecido";
  }
};

const sendEmail = async () => {
  try {
    carregandoEmails.value = true;

    const response = await axios.post(
      "http://localhost:2000/send-email",
      email.value
    );

    limpaCampos();
    form.value.resetValidation();

    setTimeout(() => {
      fetchEmails();
      carregandoEmails.value = false;
    }, 3000);
  } catch (error) {
    console.error(error);
    alert("Erro ao enviar E-mail");
  }
};

const limpaCampos = () => {
  if (form.value) {
    form.value.reset();
    form.value.resetValidation();
  }
};
</script>
