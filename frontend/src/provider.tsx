import type { NavigateOptions } from "react-router-dom";
import * as React from "react";

import { HeroUIProvider } from "@heroui/system";
import { useHref, useNavigate } from "react-router-dom";
import { Toaster } from 'react-hot-toast';
import { ThemeProvider } from '@/components/theme-provider';
import { I18nProvider } from "@react-aria/i18n";

declare module "@react-types/shared" {
  interface RouterConfig {
    routerOptions: NavigateOptions;
  }
}

export interface ProvidersProps {
  children: React.ReactNode;
}

export function Provider({ children }: ProvidersProps) {
  const navigate = useNavigate();

  return (
    <I18nProvider locale="zh-CN">
      <HeroUIProvider navigate={navigate} useHref={useHref}>
        <ThemeProvider>
          {children}
          <Toaster 
            position="top-center"
            toastOptions={{
              duration: 2000,
              className: 'dark:bg-gray-800 dark:text-white',
              style: {
                background: 'var(--toaster-bg, #ffffff)',
                color: 'var(--toaster-color, #000000)',
                border: '1px solid var(--toaster-border, #e5e7eb)',
              },
              success: {
                duration: 2000,
                style: {
                  background: '#10b981',
                  color: '#ffffff',
                },
              },
              error: {
                duration: 2000,
                style: {
                  background: '#ef4444',
                  color: '#ffffff',
                },
              },
            }}
          />
        </ThemeProvider>
      </HeroUIProvider>
    </I18nProvider>
  );
}
