
from django.contrib import admin
from django.conf.urls import url
from django.urls import path
from django.contrib.auth import views as auth_views
from django.contrib.auth.decorators import login_required
from django.urls import path, include, re_path
from .views import DashboardView

urlpatterns = [
    path('admin/', admin.site.urls),
    url(r'^login/$', auth_views.LoginView.as_view(template_name="login.html", redirect_field_name="next"), name="login"),
    path('dashboard/', login_required(DashboardView.as_view()), name="dashboard"),
    
]