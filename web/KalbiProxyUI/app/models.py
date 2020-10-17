from django.db import models

# Create your models here.


class Users(models.Model):
    id = models.AutoField(primary_key=True)
    username = models.CharField(max_length=70, null=True)
    password = models.CharField(max_length=70, null=True)
    nonce = models.CharField(max_length=70, null=True)
    